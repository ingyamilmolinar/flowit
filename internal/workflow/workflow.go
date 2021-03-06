package workflow

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/config"
)

// Workflow is the data structure representing a single workflow instance
type Workflow struct {
	ID              string
	Preffix         string
	Name            string
	IsActive        bool
	Executions      []Execution
	LatestExecution *Execution
	// TODO: Switch to pointer and refer to a separate struct in the repository
	// This to avoid several workflows based on the same config to have duplicated data in the repo
	State    config.Flowit
	Metadata WorkflowMetadata
}

// WorkflowMetadata is the data structure that provides workflow instance metadata
type WorkflowMetadata struct {
	Version  uint64
	Started  uint64
	Updated  uint64
	Finished uint64
}

// Execution is the data structure representing a single execution instance
type Execution struct {
	ID         string
	FromStage  string
	Stage      string
	Args       []string
	Checkpoint int
	Failed     bool
	Metadata   ExecutionMetadata
}

// ExecutionMetadata is the data structure that provides execution instance metadata
type ExecutionMetadata struct {
	Version  uint64
	Started  uint64
	Finished uint64
}

// OptionalWorkflow is the data type that wraps an Workflow in an optional
type OptionalWorkflow struct {
	workflow Workflow
	isSet    bool
}

// WorkflowState defines all the possible workflow states
type WorkflowState int

const (
	STARTED  WorkflowState = iota
	FAILED   WorkflowState = iota
	FINISHED WorkflowState = iota
)

// Service implements the Workflow Service methods
type Service struct{}

// NewService builds a new Workflow Service
func NewService() *Service {
	return &Service{}
}

// CreateWorkflow creates a new Workflow with a name and a variable map as inputs
func (s *Service) CreateWorkflow(workflowName string, definition config.Flowit) *Workflow {
	workflowID := uuid.New().String()
	return &Workflow{
		ID:       workflowID,
		Preffix:  workflowID[:6],
		Name:     workflowName,
		IsActive: false,
		State:    definition,
		Metadata: WorkflowMetadata{
			Version: 0,
		},
	}
}

// CancelWorkflow marks workflow as cancelled
// TODO: Unit test
func (s *Service) CancelWorkflow(w *Workflow) {
	now := uint64(time.Now().UnixNano())
	w.IsActive = false
	w.Metadata.Updated = now
	w.Metadata.Finished = now
}

// StartExecution returns a new Execution for a given Workflow
func (s *Service) StartExecution(workflow *Workflow, fromStage, currentStage string, args []string) *Execution {
	now := uint64(time.Now().UnixNano())
	execution := Execution{
		ID:         uuid.New().String(),
		FromStage:  fromStage,
		Stage:      currentStage,
		Args:       args,
		Checkpoint: -1,
		Metadata: ExecutionMetadata{
			Version: 0,
			Started: now,
		},
	}
	workflow.IsActive = true
	workflow.Executions = append([]Execution{execution}, workflow.Executions...)
	workflow.LatestExecution = &execution
	if workflow.Metadata.Started == 0 {
		workflow.Metadata.Started = now
	}
	workflow.Metadata.Updated = now
	return &execution
}

// SetCheckpoint sets the checkpoint for a given execution
func (s *Service) SetCheckpoint(execution *Execution, checkpoint int) {
	execution.Checkpoint = checkpoint
}

// FinishExecution marks a given execution as finished
func (s *Service) FinishExecution(workflow *Workflow, execution *Execution, workflowState WorkflowState) error {
	if execution.Metadata.Finished > 0 {
		return errors.New("Execution has already finished")
	}
	now := uint64(time.Now().UnixNano())
	execution.Metadata.Finished = now
	if workflowState == FAILED {
		execution.Failed = true
		execution.Stage = execution.FromStage
	}
	workflow.IsActive = workflowState != FINISHED
	workflow.Metadata.Updated = now
	if workflowState == FINISHED {
		workflow.Metadata.Finished = now
	}
	return nil
}

// AddVariables adds the given variables to the workflow instance
func (s *Service) AddVariables(workflow *Workflow, variables map[string]interface{}) {
	for k, v := range variables {
		workflow.State.Variables[k] = v
	}
}

// NewWorkflowOptional receives a Workflow and returns a filled optional ready to be unwrapped
func NewWorkflowOptional(workflow Workflow) OptionalWorkflow {
	return OptionalWorkflow{
		workflow: workflow,
		isSet:    true,
	}
}

// Get works on a pointer to OptionalWorkflow and returns the Workflow wrapped in the optional or
// returns an error in case the optional is empty
func (optional *OptionalWorkflow) Get() (Workflow, error) {
	if !optional.isSet {
		return optional.workflow, errors.New("optional value is not set")
	}
	return optional.workflow, nil
}

// Stage returns the flowit configuration stage given a stage ID
func (w Workflow) Stage(stageID string) config.Stage {
	for _, wf := range w.State.Workflows {
		if wf.ID != w.Name {
			continue
		}
		for _, s := range wf.Stages {
			if s.ID == stageID {
				return s
			}
		}
	}
	return config.Stage{}
}

// StateMachineID returns the worklow state machine ID
func (w Workflow) StateMachineID() string {
	for _, wf := range w.State.Workflows {
		if w.Name == wf.ID {
			return wf.StateMachine
		}
	}
	return ""
}
