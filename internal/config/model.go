package config

import "github.com/pkg/errors"

// WorkflowDefinition is the consumer friendly data structure that hosts the loaded workflow definition
type WorkflowDefinition struct {
	Flowit Flowit
}

// Flowit is the consumer friendly data structure that hosts the loaded workflow definition main body
type Flowit struct {
	Version       string
	Config        Config
	Variables     Variables
	StateMachines []StateMachine
	Workflows     []Workflow
}

// Config is the consumer friendly data structure that hosts the loaded workflow definition configuration
type Config struct {
	CheckpointExecution bool
	Shell               string
}

// Variables is the consumer friendly data structure that hosts the loaded workflow definition variables
type Variables map[string]interface{}

// StateMachine is the consumer friendly data structure that hosts
// the loaded workflow definition state machine
type StateMachine struct {
	ID           string
	Stages       []string
	InitialStage string
	FinalStages  []string
	Transitions  []StateMachineTransition
}

// StateMachineTransition is the consumer friendly data structure that hosts
// the loaded workflow definition state machine transition
type StateMachineTransition struct {
	From []string
	To   []string
}

// Stages is the consumer friendly data structure that hosts
// the loaded workflow definition tag stages
type Stages map[string][]string

// Workflow is the consumer friendly data structure that hosts
// the loaded workflow definition workflow
type Workflow struct {
	ID           string
	StateMachine string
	Stages       []Stage
}

// Stage is the consumer friendly data structure that hosts
// the loaded workflow definition workflow stage
type Stage struct {
	ID         string
	Args       []string
	Conditions []string
	Actions    []string
}

// Transition is the consumer friendly data structure that hosts
// the loaded workflow definition branch transition
type Transition struct {
	From string
	To   []string
}

// StateMachine receives a state machine ID and returns the correspoding
// state machine
func (wd WorkflowDefinition) StateMachine(stateMachineID string) (StateMachine, error) {
	for _, stateMachine := range wd.Flowit.StateMachines {
		if stateMachine.ID == stateMachineID {
			return stateMachine, nil
		}
	}
	return StateMachine{}, errors.New("Invalid state machine ID: " + stateMachineID)
}

// Workflow returns the workflow associated to the specified workflowID
func (wd WorkflowDefinition) Workflow(workflowID string) (Workflow, error) {
	for _, workflow := range wd.Flowit.Workflows {
		if workflow.ID == workflowID {
			return workflow, nil
		}
	}
	return Workflow{}, errors.New("Invalid workflow ID: " + workflowID)
}

// Stages returns the loaded workflow definition stages for the specified workflowID
func (wd WorkflowDefinition) Stages(workflowID string) ([]Stage, error) {
	for _, workflow := range wd.Flowit.Workflows {
		if workflow.ID == workflowID {
			return workflow.Stages, nil
		}
	}
	return nil, errors.New("Invalid workflow ID: " + workflowID)
}

// Stage returns the loaded workflow definition stage for the specified workflowID and stage
func (wd WorkflowDefinition) Stage(workflowID, stageID string) (Stage, error) {
	for _, workflow := range wd.Flowit.Workflows {
		if workflow.ID == workflowID {
			for _, stage := range workflow.Stages {
				if stage.ID == stageID {
					return stage, nil
				}
			}
			return Stage{}, errors.New("Invalid stage ID: " + stageID)
		}
	}
	return Stage{}, errors.New("Invalid workflow ID: " + workflowID)
}
