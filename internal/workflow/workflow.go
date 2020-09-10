package workflow

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Workflow is the data structure representing a single workflow instance
type Workflow struct {
	ID              string
	Preffix         string
	Name            string
	IsActive        bool
	Executions      []Execution
	LatestExecution *Execution
	Variables       map[string]interface{}
	Metadata        WorkflowMetadata
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
	ID       string
	Stage    string
	Metadata ExecutionMetadata
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

type WorkflowState int

const (
	STARTED  WorkflowState = iota
	FINISHED WorkflowState = iota
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// CreateWorkflow creates a new Workflow with a name and a variable map as inputs
func (s *Service) CreateWorkflow(workflowName string, variables map[string]interface{}) *Workflow {
	workflowID := uuid.New().String()

	return &Workflow{
		ID:        workflowID,
		Preffix:   workflowID[:6],
		Name:      workflowName,
		IsActive:  false,
		Variables: variables,
		Metadata: WorkflowMetadata{
			Version: 0,
		},
	}
}

func (s *Service) StartExecution(workflow *Workflow, currentStage string) *Execution {
	now := uint64(time.Now().UnixNano())
	execution := Execution{
		ID:    uuid.New().String(),
		Stage: currentStage,
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

func (s *Service) FinishExecution(workflow *Workflow, execution *Execution, workflowState WorkflowState) error {
	if execution.Metadata.Finished > 0 {
		return errors.New("Execution has already finished")
	}
	now := uint64(time.Now().UnixNano())
	execution.Metadata.Finished = now
	workflow.IsActive = workflowState != FINISHED
	workflow.Metadata.Updated = now
	if workflowState == FINISHED {
		workflow.Metadata.Finished = now
	}
	return nil
}

func (s *Service) AddVariables(workflow *Workflow, variables map[string]interface{}) {
	for k, v := range variables {
		workflow.Variables[k] = v
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
