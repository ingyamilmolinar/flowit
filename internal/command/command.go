package command

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/io"
	"github.com/yamil-rivera/flowit/internal/utils"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

type RuntimeService interface {
	Execute(commands []string, variables map[string]interface{}, checkpoint int) ([]string, int, error)
}

type FsmService interface {
	OriginState() string
	InitialState(stateMachineID string) string
	AvailableStates(stateMachineID string, state string) []string
	IsFinalState(stateMachineID, state string) bool
}

type RepositoryService interface {
	GetWorkflow(workflowName, workflowID string) (w.OptionalWorkflow, error)
	GetWorkflows(workflowName string, count int, excludeInactive bool) ([]w.Workflow, error)
	GetWorkflowFromPreffix(workflowName, workflowIDPreffix string) (w.OptionalWorkflow, error)
	PutWorkflow(workflow w.Workflow) error
	DeleteWorkflow(workflowName, workflowID string) error
}

type WorkflowService interface {
	CreateWorkflow(workflowName string, variables map[string]interface{}) *w.Workflow
	CancelWorkflow(workflow *w.Workflow)
	StartExecution(workflow *w.Workflow, fromStage, currentState string) *w.Execution
	SetCheckpoint(execution *w.Execution, checkpoint int)
	FinishExecution(workflow *w.Workflow, execution *w.Execution, workflowState w.WorkflowState) error
	AddVariables(workflow *w.Workflow, variables map[string]interface{})
}

type Service struct {
	rootCommand        *cobra.Command
	runtimeService     RuntimeService
	fsmService         FsmService
	repositoryService  RepositoryService
	workflowService    WorkflowService
	workflowDefinition *config.WorkflowDefinition
}

type command struct {
	cobra       *cobra.Command
	subcommands []command
}

func NewService(run RuntimeService, fsm FsmService, repo RepositoryService, wf WorkflowService, wd *config.WorkflowDefinition) *Service {
	return &Service{nil, run, fsm, repo, wf, wd}
}

// RegisterCommands registers all commands and subcommands based on the provided configuration
func (s *Service) RegisterCommands(version string) error {

	workflowDefinitions := s.workflowDefinition.Flowit.Workflows
	mainCommands := make([]command, len(workflowDefinitions)+1)
	for i, workflowDefinition := range workflowDefinitions {

		workflowName := workflowDefinition.ID
		mainCommands[i].cobra = newContainerCommand(workflowName)
		activeWorkflows, err := s.getAllActiveWorkflows(workflowName)

		if err != nil {
			return errors.WithStack(err)
		}

		if len(activeWorkflows) > 0 {

			mainCommands[i].subcommands = make([]command, len(activeWorkflows))
			for j, activeWorkflow := range activeWorkflows {

				mainCommands[i].subcommands[j].cobra = newContainerCommand(activeWorkflow.Preffix)
				stages, err := s.generatePossibleCommands(workflowName, activeWorkflow.ID)
				if err != nil {
					return errors.Wrap(err, "Error generating possible commands")
				}

				mainCommands[i].subcommands[j].subcommands = stages
			}

			initialStages, err := s.generateInitialCommands(workflowName)
			if err != nil {
				return errors.WithStack(err)
			}

			mainCommands[i].subcommands = append(mainCommands[i].subcommands, initialStages...)
		} else {

			initialStages, err := s.generateInitialCommands(workflowName)
			if err != nil {
				return errors.WithStack(err)
			}

			mainCommands[i].subcommands = initialStages
		}
	}

	mainCommands[len(workflowDefinitions)].cobra = newSimpleCommand("version", version)

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

func (s Service) getAllActiveWorkflows(workflowName string) ([]w.Workflow, error) {
	return s.repositoryService.GetWorkflows(workflowName, 0, true)
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

func newStageCommand(use string, args int, run func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:  use,
		Args: cobra.ExactArgs(args),
		RunE: run,
	}
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
				var workflow *w.Workflow
				if !optionalWorkflowID.IsSet() {
					workflow = s.workflowService.CreateWorkflow(workflowName, s.workflowDefinition.Flowit.Variables)
					io.Println("Workflow with ID: " + workflow.ID + " was created")
				} else {
					workflowID, _ := optionalWorkflowID.Get()
					optionalWorkflow, _ := s.repositoryService.GetWorkflowFromPreffix(workflowName, workflowID)
					wf, _ := optionalWorkflow.Get()
					workflow = &wf
				}

				checkpoint := 0
				// TODO: This should only apply if the arguments are the same
				if workflow.LatestExecution != nil && s.workflowDefinition.Flowit.Config.CheckpointExecution {
					lastExecution := workflow.LatestExecution
					if lastExecution.Checkpoint >= 0 {
						checkpoint = lastExecution.Checkpoint
					}
				}

				fromStageID := s.fsmService.OriginState()
				if workflow.LatestExecution != nil {
					fromStageID = workflow.LatestExecution.Stage
				}

				execution := s.workflowService.StartExecution(workflow, fromStageID, stageID)
				stage, _ := s.workflowDefinition.Stage(workflowName, stageID)

				if len(stage.Args) > 0 {
					variables := make(map[string]interface{})
					for i, arg := range stage.Args {
						variable, err := utils.ExtractVariableNameFromVariableDeclaration(arg)
						if err != nil {
							return errors.WithStack(err)
						}
						variables[variable] = args[i]
					}
					s.workflowService.AddVariables(workflow, variables)
				}

				if len(stage.Conditions) > 0 {
					io.Println("Running conditions...")
					out, _, err := s.runtimeService.Execute(stage.Conditions, workflow.Variables, 0)
					io.Println(strings.Join(utils.MergeSlices(stage.Conditions, out), "\n"))
					if err != nil {
						return errors.WithStack(err)
					}
				}

				io.Println("Running actions...")
				out, failedActionIdx, err := s.runtimeService.Execute(stage.Actions, workflow.Variables, checkpoint)
				if err != nil {
					io.Println(strings.Join(utils.MergeSlices(stage.Actions[checkpoint:failedActionIdx], out), "\n"))
					if s.workflowDefinition.Flowit.Config.CheckpointExecution {
						s.workflowService.SetCheckpoint(execution, failedActionIdx)
						io.Println("Checkpoint set on command: ", stage.Actions[failedActionIdx])
						execution.Stage = execution.FromStage
						if err := s.workflowService.FinishExecution(workflow, execution, w.FAILED); err != nil {
							return errors.WithStack(err)
						}
						if err := s.repositoryService.PutWorkflow(*workflow); err != nil {
							return errors.WithStack(err)
						}
					}
					return errors.WithStack(err)
				}
				io.Println(strings.Join(utils.MergeSlices(stage.Actions[checkpoint:], out), "\n"))

				isFinal := s.fsmService.IsFinalState(workflowName, stageID)
				workflowState := w.STARTED
				if isFinal {
					workflowState = w.FINISHED
				}
				err = s.workflowService.FinishExecution(workflow, execution, workflowState)
				if err != nil {
					return errors.WithStack(err)
				}
				if err := s.repositoryService.PutWorkflow(*workflow); err != nil {
					return errors.WithStack(err)
				}
				return nil

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

func (s Service) generateInitialCommands(workflowName string) ([]command, error) {

	initialEvent := s.fsmService.InitialState(workflowName)
	initialEvents := []string{initialEvent}
	commands, err := s.generateCommandsFromStages(workflowName, initialEvents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return commands, nil
}

func (s Service) generatePossibleCommands(workflowName, workflowID string) ([]command, error) {

	workflowOptional, err := s.repositoryService.GetWorkflow(workflowName, workflowID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var availableStates []string
	workflow, err := workflowOptional.Get()
	if err == nil {
		if workflow.LatestExecution.Checkpoint >= 0 {
			availableStates = s.fsmService.AvailableStates(workflowName, workflow.LatestExecution.FromStage)
		} else {
			availableStates = s.fsmService.AvailableStates(workflowName, workflow.LatestExecution.Stage)
		}
	} else {
		availableState := s.fsmService.InitialState(workflowName)
		availableStates = append(availableStates, availableState)
	}

	commands, err := s.generateCommandsFromStages(workflowName, availableStates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commands = append(commands, s.generateCancelCommand(workflowName))
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

					// We are sure optionalWorkflowID is always wrapping a workflowID
					workflowID, _ := optionalWorkflowID.Get()
					workflowOptional, err := s.repositoryService.GetWorkflow(workflowName, workflowID)
					if err != nil {
						return errors.WithStack(err)
					}
					workflow, err := workflowOptional.Get()
					if err != nil {
						return errors.WithStack(err)
					}
					s.workflowService.CancelWorkflow(&workflow)
					if err := s.repositoryService.PutWorkflow(workflow); err != nil {
						return errors.WithStack(err)
					}
					return io.Println("Workflow with ID: " + workflowID + " was cancelled")
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
