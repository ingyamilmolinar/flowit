package config

import (
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

// ProcessWorkflowDefinition reads, parses the specified yaml configuration file and returns a map with the key/values
func ProcessWorkflowDefinition(fileName string, fileLocation string) (*WorkflowDefinition, error) {

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

	// Since viper does not allow for array defaults, we roll our own mechanism
	setDefaults(rawWorkflowDefinition)

	var workflowDefinition WorkflowDefinition
	if err := utils.DeepCopy(rawWorkflowDefinition, &workflowDefinition); err != nil {
		return nil, errors.WithStack(err)
	}
	return &workflowDefinition, nil
}
