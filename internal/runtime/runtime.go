package runtime

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/fsm"
	"github.com/yamil-rivera/flowit/internal/utils"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

// Service exposes the methods to interact with the Runtime Service
type Service struct {
	repositoryService RepositoryService
	fsmServiceFactory fsm.FsmServiceFactory
	workflowService   WorkflowService
}

// RepositoryService defines the methods that must be implemented in order for a struct to be considered a Repository Service by the RuntimeService
type RepositoryService interface {
	GetWorkflow(workflowName, workflowID string) (w.OptionalWorkflow, error)
	GetWorkflows(workflowName string, count int, excludeInactive bool) ([]w.Workflow, error)
	GetAllWorkflows(excludeInactive bool) ([]w.Workflow, error)
	GetWorkflowFromPreffix(workflowName, workflowIDPreffix string) (w.OptionalWorkflow, error)
	PutWorkflow(workflow w.Workflow) error
}

// WorkflowService defines the methods that must be implemented in order for a struct to be considered a Workflow Service by the RuntimeService
type WorkflowService interface {
	CreateWorkflow(workflowName string, definition config.Flowit) *w.Workflow
	CancelWorkflow(workflow *w.Workflow)
	StartExecution(workflow *w.Workflow, fromStage, currentState string, args []string) *w.Execution
	SetCheckpoint(execution *w.Execution, checkpoint int)
	FinishExecution(workflow *w.Workflow, execution *w.Execution, workflowState w.WorkflowState) error
	AddVariables(workflow *w.Workflow, variables map[string]interface{})
}

// Writer defines the methods that must be implemented in order for a struct to be considered a Writer by the RuntimeService
// A Writer is an object which encapsulates a write side-effect
// It is used by the RuntimeService to avoid depending on a concrete logging implementation
type Writer interface {
	Write(s string) error
}

// Executor defines the methods that must be implemented in order for a struct to be considered an Executor by the RuntimeService
type Executor interface {
	Config(shell string)
	Execute(command string) (string, error)
}

// UnixShellExecutor is the default implementation of the Executor interface
type UnixShellExecutor struct {
	shell string
}

// NewService returns a new instance of the RuntimeService
func NewService(rs RepositoryService, fsf fsm.FsmServiceFactory, ws WorkflowService) *Service {
	return &Service{rs, fsf, ws}
}

// NewUnixShellExecutor returns an Executor instance based on the UnixShellExecutor
func NewUnixShellExecutor() Executor {
	return &UnixShellExecutor{}
}

// Config configures the UnixShellExecutor using a shell binary location
func (e *UnixShellExecutor) Config(shell string) {
	e.shell = shell
}

// TODO: Handle && exit 1
// Execute receives a command, runs it using the configured shell and returns the produced output
func (e *UnixShellExecutor) Execute(command string) (string, error) {
	shellArgs := strings.Split(e.shell, " ")
	mainCommand := shellArgs[0]
	restOfArgs := append(shellArgs[1:], "-c", command)
	cmd := exec.Command(mainCommand, restOfArgs...)
	out, err := cmd.Output()
	trimmedOut := strings.TrimSuffix(string(out), "\n")
	if err != nil {
		return trimmedOut, errors.Wrap(err, "Error executing command: "+command+" with shell: "+e.shell)
	}
	return trimmedOut, nil
}

// Run executes a workflow stage based on the provided configuration or based on a persisted workflow
// If optionalWorkflowPreffix is not empty, the workflow state will be retrieved from the repository
// If optionalWorkflowPreffix is empty, the provided workflow definition will be used to create a new workflow in the repository
func (s *Service) Run(optionalWorkflowPreffix utils.OptionalString, args []string, workflowName, stageID string, workflowDefinition config.Flowit, executor Executor, writer Writer) error {
	var workflow *w.Workflow
	if !optionalWorkflowPreffix.IsSet() {
		workflow = s.workflowService.CreateWorkflow(workflowName, workflowDefinition)
		// nolint: errcheck
		writer.Write("Workflow with ID: " + workflow.ID + " was created")
	} else {
		workflowPreffix, _ := optionalWorkflowPreffix.Get()
		optionalWorkflow, _ := s.repositoryService.GetWorkflowFromPreffix(workflowName, workflowPreffix)
		wf, _ := optionalWorkflow.Get()
		workflow = &wf
	}
	fsmService, err := s.fsmServiceFactory.NewFsmService(workflow.State)
	if err != nil {
		return errors.WithStack(err)
	}

	fromStageID := fsmService.OriginState()
	if workflow.LatestExecution != nil {
		fromStageID = workflow.LatestExecution.Stage
	}

	stage := workflow.Stage(stageID)
	if !fsmService.IsTransitionValid(workflow.StateMachineID(), fromStageID, stage.ID) {
		return errors.Errorf("Invalid transition from %s to %s", fromStageID, stageID)
	}

	checkpoint := 0
	if workflow.LatestExecution != nil && workflow.State.Config.CheckpointExecution {
		lastExecution := workflow.LatestExecution
		if lastExecution.Failed && !utils.CompareSlices(lastExecution.Args, args) {
			return errors.Errorf("Arguments: %+v do not match with last failed execution arguments: %+v", args, lastExecution.Args)
		}
		if lastExecution.Checkpoint >= 0 {
			checkpoint = lastExecution.Checkpoint
		}
	}

	execution := s.workflowService.StartExecution(workflow, fromStageID, stageID, args)

	if len(stage.Args) > 0 {
		if len(args) != len(stage.Args) {
			return errors.Errorf("Wrong number of arguments provided. Expected %d but got %d.", len(stage.Args), len(args))
		}
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

	// Set executor for this run based on workflow state
	executor.Config(workflow.State.Config.Shell)

	err = s.runConditions(stage.Conditions, workflow.State.Variables, executor, writer)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.runActions(workflow, execution, stage.Actions, workflow.State.Variables, workflow.State.Config.CheckpointExecution, checkpoint, executor, writer)
	if err != nil {
		return errors.WithStack(err)
	}

	stateMachineID := workflow.StateMachineID()
	isFinal := fsmService.IsFinalState(stateMachineID, stageID)
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

// Cancel marks the provided workflowID as cancelled
func (s *Service) Cancel(workflowID string, workflowName string, writer Writer) error {
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
	// nolint: errcheck
	writer.Write("Workflow with ID: " + workflow.ID + " was cancelled")
	return nil
}

func (s Service) execute(commands []string, variables map[string]interface{}, checkpoint int, executor Executor, writer Writer) (int, error) {

	i, err := s.runCommands(commands[checkpoint:], variables, executor, writer)
	if err != nil {
		return i + checkpoint, errors.WithStack(err)
	}
	return 0, nil
}

func (s Service) runCommands(commands []string, variables map[string]interface{}, executor Executor, writer Writer) (int, error) {

	if len(commands) == 0 {
		return 0, nil
	}

	for i, command := range commands {
		parsedCommand, err := utils.EvaluateVariablesInExpression(command, variables)
		if err != nil {
			return i, errors.Wrap(err, "Error evaluating variables in command: "+command)
		}
		out, err := executor.Execute(parsedCommand)
		// nolint: errcheck
		writer.Write(out)
		if err != nil {
			return i, errors.Wrap(err, "Error executing command: "+parsedCommand)
		}
	}
	return 0, nil
}

func (s Service) runConditions(conditions []string, variables map[string]interface{}, executor Executor, writer Writer) error {
	if len(conditions) > 0 {
		// nolint: errcheck
		writer.Write("Running conditions...")
		_, err := s.execute(conditions, variables, 0, executor, writer)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s Service) runActions(workflow *w.Workflow, execution *w.Execution, actions []string, variables map[string]interface{}, checkpointEnabled bool, checkpoint int, executor Executor, writer Writer) error {
	// nolint: errcheck
	writer.Write("Running actions...")
	failedActionIdx, err := s.execute(actions, variables, checkpoint, executor, writer)
	if err != nil {
		// TOFIX:
		// stdout = append(stdout, utils.MergeSlices(actions[checkpoint:failedActionIdx], out)...)
		if checkpointEnabled {
			s.workflowService.SetCheckpoint(execution, failedActionIdx)
			// nolint: errcheck
			writer.Write("Checkpoint set on command: " + actions[failedActionIdx])
			if err := s.workflowService.FinishExecution(workflow, execution, w.FAILED); err != nil {
				return errors.WithStack(err)
			}
			if err := s.repositoryService.PutWorkflow(*workflow); err != nil {
				return errors.WithStack(err)
			}
		}
		return errors.WithStack(err)
	}
	// TOFIX:
	// stdout = append(stdout, utils.MergeSlices(actions[checkpoint:], out)...)
	return nil
}
