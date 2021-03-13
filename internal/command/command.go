package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/fsm"
	"github.com/yamil-rivera/flowit/internal/io"
	"github.com/yamil-rivera/flowit/internal/runtime"
	"github.com/yamil-rivera/flowit/internal/utils"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

type RuntimeService interface {
	Run(optionalWorkflowID utils.OptionalString, args []string, workflowName, stageID string, workflowDefinition config.Flowit, executor runtime.Executor, writer runtime.Writer) error
	Cancel(optionalWorkflowID utils.OptionalString, args []string, workflowName string, writer runtime.Writer) error
}

type RepositoryService interface {
	GetWorkflow(workflowName, workflowID string) (w.OptionalWorkflow, error)
	GetWorkflows(workflowName string, count int, excludeInactive bool) ([]w.Workflow, error)
	GetAllWorkflows(excludeInactive bool) ([]w.Workflow, error)
	GetWorkflowFromPreffix(workflowName, workflowIDPreffix string) (w.OptionalWorkflow, error)
	PutWorkflow(workflow w.Workflow) error
}

type Service struct {
	rootCommand        *cobra.Command
	runtimeService     RuntimeService
	fsmServiceFactory  fsm.FsmServiceFactory
	repositoryService  RepositoryService
	workflowDefinition *config.WorkflowDefinition
}

type command struct {
	cobra       *cobra.Command
	subcommands []command
}

func NewService(run RuntimeService, fsf fsm.FsmServiceFactory, repo RepositoryService, wd *config.WorkflowDefinition) *Service {
	return &Service{nil, run, fsf, repo, wd}
}

// RegisterCommands registers all commands and subcommands based on the provided configuration
// and previous active workflows
func (s *Service) RegisterCommands(version string) error {

	var mainCommands []command
	fsmService, err := s.fsmServiceFactory.NewFsmService(s.workflowDefinition.Flowit)
	if err != nil {
		return errors.WithStack(err)
	}
	workflowDefinitions := s.workflowDefinition.Flowit.Workflows
	for _, workflowDefinition := range workflowDefinitions {

		workflowName := workflowDefinition.ID
		stateMachine := workflowDefinition.StateMachine
		cmd := command{}
		cmd.cobra = newContainerCommand(workflowName)
		initialStages, err := s.generateInitialCommands(fsmService, stateMachine, workflowName)
		if err != nil {
			return errors.WithStack(err)
		}
		cmd.subcommands = initialStages
		mainCommands = append(mainCommands, cmd)
	}

	activeWorkflows, err := s.getAllActiveWorkflows()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, workflow := range activeWorkflows {
		childCmd := command{}
		childCmd.cobra = newContainerCommand(workflow.Preffix)
		stages, err := s.generatePossibleCommands(workflow)
		if err != nil {
			return errors.Wrap(err, "Error generating possible commands")
		}
		childCmd.subcommands = stages

		// Check if we already have a registered command for this workflow name
		var cmd *command
		var found bool
		for _, mainCmd := range mainCommands {
			if mainCmd.cobra.Use == workflow.Name {
				cmd = &mainCmd
				found = true
			}
		}

		if !found {
			cmd = &command{}
			cmd.cobra = newContainerCommand(workflow.Name)
		}
		cmd.subcommands = append(cmd.subcommands, childCmd)
		mainCommands = replaceCommand(mainCommands, *cmd)
	}

	cmd := command{}
	cmd.cobra = newSimpleCommand("version", version)
	mainCommands = append(mainCommands, cmd)

	rootCommand := &cobra.Command{
		Use:   "flowit",
		Short: "A flexible workflow manager",
		Long:  "A flexible workflow manager",
	}

	for _, mainCommand := range mainCommands {
		for _, subcommands := range mainCommand.subcommands {
			for _, subcommand := range subcommands.subcommands {
				subcommands.cobra.AddCommand(subcommand.cobra)
			}
			mainCommand.cobra.AddCommand(subcommands.cobra)
		}
		rootCommand.AddCommand(mainCommand.cobra)
	}

	s.rootCommand = rootCommand
	return nil
}

func (s Service) Execute() error {
	if err := s.rootCommand.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s Service) getAllActiveWorkflowsByName(workflowName string) ([]w.Workflow, error) {
	return s.repositoryService.GetWorkflows(workflowName, 0, true)
}

func (s Service) getAllActiveWorkflows() ([]w.Workflow, error) {
	return s.repositoryService.GetAllWorkflows(true)
}

func newContainerCommand(commandUse string) *cobra.Command {
	return &cobra.Command{
		Use: commandUse,
	}
}

func newSimpleCommand(commandUse string, commandRun string) *cobra.Command {
	return &cobra.Command{
		Use: commandUse,
		Run: func(cmd *cobra.Command, args []string) {
			io.Println(commandRun)
		},
	}
}

// TODO: Add arguments description to command help
func newStageCommand(use string, args int, run func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:  use,
		Args: cobra.ExactArgs(args),
		RunE: run,
	}
}

func (s Service) generateCommandsFromStagesForWorkflow(workflow w.Workflow, stages []string) ([]command, error) {
	commands := make([]command, len(stages))
	for i, stageID := range stages {

		runFunc := func(workflowName string, stageID string) func(cmd *cobra.Command, args []string) error {

			return func(cmd *cobra.Command, args []string) error {
				optionalWorkflowID, err := s.getWorkflowIDFromCommand(cmd)
				if err != nil {
					return errors.WithStack(err)
				}
				err = s.runtimeService.Run(optionalWorkflowID, args, workflowName, stageID, s.workflowDefinition.Flowit, runtime.NewUnixShellExecutor(), io.NewConsoleWriter())
				return err
			}

		}(workflow.Name, stageID)

		stage, err := stage(workflow, workflow.Name, stageID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		commands[i].cobra = newStageCommand(stage.ID, len(stage.Args), runFunc)

	}
	return commands, nil
}

func (s Service) generateCommandsFromStages(workflowName string, stages []string) ([]command, error) {
	commands := make([]command, len(stages))
	for i, stageID := range stages {

		runFunc := func(workflowName string, stageID string) func(cmd *cobra.Command, args []string) error {

			return func(cmd *cobra.Command, args []string) error {
				optionalWorkflowID, err := s.getWorkflowIDFromCommand(cmd)
				if err != nil {
					return errors.WithStack(err)
				}
				err = s.runtimeService.Run(optionalWorkflowID, args, workflowName, stageID, s.workflowDefinition.Flowit, runtime.NewUnixShellExecutor(), io.NewConsoleWriter())
				return err
			}

		}(workflowName, stageID)

		stage, err := s.workflowDefinition.Stage(workflowName, stageID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		commands[i].cobra = newStageCommand(stage.ID, len(stage.Args), runFunc)

	}
	return commands, nil
}

func (s Service) generateInitialCommands(fsmService fsm.Service, stateMachine, workflowName string) ([]command, error) {

	initialEvent := fsmService.InitialState(stateMachine)
	initialEvents := []string{initialEvent}
	commands, err := s.generateCommandsFromStages(workflowName, initialEvents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return commands, nil
}

func (s Service) generatePossibleCommands(workflow w.Workflow) ([]command, error) {
	fsmService, err := s.fsmServiceFactory.NewFsmService(workflow.State)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var availableStates []string
	if workflow.LatestExecution.Checkpoint >= 0 {
		availableStates = fsmService.AvailableStates(workflow.StateMachineID(), workflow.LatestExecution.FromStage)
	} else {
		availableStates = fsmService.AvailableStates(workflow.StateMachineID(), workflow.LatestExecution.Stage)
	}

	commands, err := s.generateCommandsFromStagesForWorkflow(workflow, availableStates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commands = append(commands, s.generateCancelCommand(workflow.Name))
	return commands, nil

}

func (s Service) generateCancelCommand(workflowName string) command {

	return command{
		cobra: &cobra.Command{
			Use: "cancel",
			RunE: func(workflowName string) func(cmd *cobra.Command, args []string) error {

				return func(cmd *cobra.Command, args []string) error {
					optionalWorkflowID, err := s.getWorkflowIDFromCommand(cmd)
					if err != nil {
						return errors.WithStack(err)
					}
					err = s.runtimeService.Cancel(optionalWorkflowID, args, workflowName, io.NewConsoleWriter())
					return err
				}

			}(workflowName),
		},
	}

}

// cmd parent is either a workflow definition name or a workflow instance name
func (s Service) getWorkflowIDFromCommand(cmd *cobra.Command) (utils.OptionalString, error) {

	parentName := cmd.Parent().Name()
	isWorkflowDefinition := false
	for _, workflow := range s.workflowDefinition.Flowit.Workflows {
		if workflow.ID == parentName {
			isWorkflowDefinition = true
			break
		}
	}

	if !isWorkflowDefinition {
		// parentName has to be a workflow instance name and cmd.Parent().Parent() a workflow definition
		workflowID, err := s.getWorkflowIDFromName(cmd.Parent().Parent().Name(), parentName)
		if err == nil {
			return utils.NewStringOptional(workflowID), nil
		}
		return utils.OptionalString{}, errors.WithStack(err)
	}

	return utils.OptionalString{}, nil
}

func (s Service) getWorkflowIDFromName(workflowName, workflowPreffix string) (string, error) {

	optionalWorkflow, err := s.repositoryService.GetWorkflowFromPreffix(workflowName, workflowPreffix)
	if err != nil {
		return "", errors.WithStack(err)
	}

	workflow, err := optionalWorkflow.Get()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return workflow.ID, nil
}

func replaceCommand(cmds []command, cmd command) []command {
	var result []command
	for _, c := range cmds {
		if c.cobra.Use == cmd.cobra.Use {
			result = append(result, cmd)
			continue
		}
		result = append(result, c)
	}
	return result
}

func stringifyCmd(c *command) string {
	result := c.cobra.Name() + " "
	for i, s := range c.subcommands {
		result += fmt.Sprint(i) + ":" + stringifyCmd(&s)
	}
	return result
}

func stage(w w.Workflow, workflowID, stageID string) (config.Stage, error) {
	for _, workflow := range w.State.Workflows {
		if workflow.ID == workflowID {
			for _, stage := range workflow.Stages {
				if stage.ID == stageID {
					return stage, nil
				}
			}
			return config.Stage{}, errors.New("Invalid stage ID: " + stageID)
		}
	}
	return config.Stage{}, errors.New("Invalid workflow ID: " + workflowID)
}
