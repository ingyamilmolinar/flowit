package models

import "errors"

// Workflow is the data structure representing a single workflow instance.
// It's life expands every state change (execution) since the initial stage up to and including any of the final stages.
type Workflow struct {
	ID              string
	Name            string
	DefinitionID    string
	IsActive        bool
	Executions      []Execution
	LatestExecution *Execution
	Variables       map[string]string
	Metadata        WorkflowMetadata
}

// WorkflowMetadata is the data structure that provides workflow instance metadata such as when it started,
// when it was most recently updated and when it ended
type WorkflowMetadata struct {
	Version  uint64
	Started  uint64
	Updated  uint64
	Finished uint64
}

// Execution is the data structure that represents a single instance every time a state change is performed
type Execution struct {
	ID       string
	State    string
	Metadata ExecutionMetadata
}

// ExecutionMetadata is the data structure that provides single state change metadata such as when it started
// and when it ended
type ExecutionMetadata struct {
	Version  uint64
	Started  uint64
	Finished uint64
}

// OptionalExecution is the data type that wraps an Execution in an optional
type OptionalExecution struct {
	execution Execution
	isSet     bool
}

// NewExecution receives an Execution and returns a filled optional ready to be unwrapped
func NewExecution(exec Execution) OptionalExecution {
	return OptionalExecution{
		execution: exec,
		isSet:     true,
	}
}

// Get works on a pointer to OptionalExecution and returns the Execution wrapped in the optional or
// returns an error in case the optional is empty
func (optional *OptionalExecution) Get() (Execution, error) {
	if !optional.isSet {
		return optional.execution, errors.New("optional value is not set")
	}
	return optional.execution, nil
}

// OptionalWorkflow is the data type that wraps an Workflow in an optional
type OptionalWorkflow struct {
	workflow Workflow
	isSet    bool
}

// NewWorkflow receives a Workflow and returns a filled optional ready to be unwrapped
func NewWorkflow(workflow Workflow) OptionalWorkflow {
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
