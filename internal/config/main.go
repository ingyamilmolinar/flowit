package config

import (
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

type ConfigService struct {
	workflowDefinition *WorkflowDefinition
}

// NewConfigService reads, parses and validates the specified configuration file and creates a new config service
// Returns an error if any step fails
func NewConfigService(fileName string, fileLocation string) (*ConfigService, error) {
	// TODO: Hash parsed and validated config and verify if it changed or not?
	viper, err := readWorkflowDefinition(fileName, fileLocation)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: Viper is allowing repeated keys...
	rawWorkflowDefinition, err := unmarshallWorkflowDefinition(viper)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = validateWorkflowDefinition(rawWorkflowDefinition); err != nil {
		return nil, errors.WithStack(err)
	}

	applyTransformations(rawWorkflowDefinition)

	// Since viper does not allow for array defaults, we roll our own mechanism
	setDefaults(rawWorkflowDefinition)

	var workflowDefinition WorkflowDefinition
	if err := utils.DeepCopy(rawWorkflowDefinition, &workflowDefinition); err != nil {
		return nil, errors.WithStack(err)
	}
	return &ConfigService{&workflowDefinition}, nil
}

// GetVersion returns the loaded workflow definition version
func (cs ConfigService) GetVersion() string {
	return cs.workflowDefinition.Flowit.Version
}

// GetConfig returns the loaded workflow definition configuration
func (cs ConfigService) GetConfig() Config {
	return cs.workflowDefinition.Flowit.Config
}

// GetVariables returns the loaded workflow definition variables
func (cs ConfigService) GetVariables() Variables {
	return cs.workflowDefinition.Flowit.Variables
}

// GetBranches returns the loaded workflow definition branches
func (cs ConfigService) GetBranches() []Branch {
	return cs.workflowDefinition.Flowit.Branches
}

// GetTags returns the loaded workflow definition tags
func (cs ConfigService) GetTags() []Tag {
	return cs.workflowDefinition.Flowit.Tags
}

// GetStateMachines returns the loaded workflow definition state machines
func (cs ConfigService) GetStateMachines() []StateMachine {
	return cs.workflowDefinition.Flowit.StateMachines
}

// GetWorkflows returns the loaded workflow definition workflows
func (cs ConfigService) GetWorkflows() []Workflow {
	return cs.workflowDefinition.Flowit.Workflows
}

// GetStages returns the loaded workflow definition stages for the specified workflowID
func (cs ConfigService) GetStages(workflowID string) ([]Stage, error) {
	for _, workflow := range cs.workflowDefinition.Flowit.Workflows {
		if workflow.ID == workflowID {
			return workflow.Stages, nil
		}
	}
	return nil, errors.New("Invalid workflowID: " + workflowID)
}

// GetStage returns the loaded workflow definition stage for the specified workflowID and stageID
func (cs ConfigService) GetStage(workflowID, stageID string) (Stage, error) {
	for _, workflow := range cs.workflowDefinition.Flowit.Workflows {
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
