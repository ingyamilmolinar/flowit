package config

import (
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var workflowDefinition WorkflowDefinition
var configLoaded bool

// LoadConfiguration reads, parses and validates the specified yaml configuration file
// Returns an error if any step fails
func LoadConfiguration(fileName string, fileLocation string) error {

	// TODO: Hash parsed and validated config and verify if it changed or not?
	viper, err := readWorkflowDefinition(fileName, fileLocation)
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: Viper is allowing repeated keys...
	rawWorkflowDefinition, err := unmarshallWorkflowDefinition(viper)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = validateWorkflowDefinition(rawWorkflowDefinition); err != nil {
		return errors.WithStack(err)
	}

	applyTransformations(rawWorkflowDefinition)

	// Since viper does not allow for array defaults, we roll our own mechanism
	setDefaults(rawWorkflowDefinition)

	if err := utils.DeepCopy(rawWorkflowDefinition, &workflowDefinition); err != nil {
		return errors.WithStack(err)
	}
	configLoaded = true
	return nil
}

// GetVersion returns the loaded workflow definition version
// Panics if configuration has not been loaded
func GetVersion() string {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Version
}

// GetConfig returns the loaded workflow definition configuration
// Panics if configuration has not been loaded
func GetConfig() Config {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Config
}

// GetVariables returns the loaded workflow definition variables
// Panics if configuration has not been loaded
func GetVariables() Variables {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Variables
}

// GetBranches returns the loaded workflow definition branches
// Panics if configuration has not been loaded
func GetBranches() []Branch {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Branches
}

// GetTags returns the loaded workflow definition tags
// Panics if configuration has not been loaded
func GetTags() []Tag {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Tags
}

// GetStateMachines returns the loaded workflow definition state machines
// Panics if configuration has not been loaded
func GetStateMachines() []StateMachine {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.StateMachines
}

// GetWorkflows returns the loaded workflow definition workflows
// Panics if configuration has not been loaded
func GetWorkflows() []Workflow {
	panicIfConfigNotLoaded()
	return workflowDefinition.Flowit.Workflows
}

// GetStages returns the loaded workflow definition stages for the specified workflowID
// Panics if configuration has not been loaded
func GetStages(workflowID string) ([]Stage, error) {
	panicIfConfigNotLoaded()
	for _, workflow := range workflowDefinition.Flowit.Workflows {
		if workflow.ID == workflowID {
			return workflow.Stages, nil
		}
	}
	return nil, errors.New("Invalid workflowID: " + workflowID)
}

// GetStage returns the loaded workflow definition stage for the specified workflowID and stageID
// Panics if configuration has not been loaded
func GetStage(workflowID, stageID string) (Stage, error) {
	panicIfConfigNotLoaded()
	for _, workflow := range workflowDefinition.Flowit.Workflows {
		if workflow.ID == workflowID {
			for _, stage := range workflow.Stages {
				if stage.ID == stageID {
					return stage, nil
				}
			}
			return Stage{}, errors.New("Invalid stageID: " + stageID)
		}
	}
	return Stage{}, errors.New("Invalid workflowID: " + workflowID)
}

func panicIfConfigNotLoaded() {
	if !configLoaded {
		panic("Configuration has not been loaded successfully!")
	}
}
