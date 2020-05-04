package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func unmarshallWorkflowDefinition(v *viper.Viper) (*rawWorkflowDefinition, error) {

	var workflowDefinition rawWorkflowDefinition

	config := func(c *mapstructure.DecoderConfig) {
		c.ErrorUnused = true
		c.WeaklyTypedInput = false
		c.ZeroFields = true
	}

	if err := (*v).UnmarshalExact(&workflowDefinition, config); err != nil {
		return nil, errors.Wrap(err, "Workflow definition unmarshalling error")
	}

	return &workflowDefinition, nil
}
