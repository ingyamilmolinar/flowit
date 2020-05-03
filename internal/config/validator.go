package config

import (
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

// ValidateConfig takes a raw configuration and validates it section by section
func validateConfig(workflowDefinition *rawWorkflowDefinition) error {
	if workflowDefinition.Flowit == nil {
		return errors.New("Configuration must have a 'flowit' key")
	}
	return validator.ValidateStruct(workflowDefinition.Flowit, configFieldRules(workflowDefinition.Flowit)...)
}

func configFieldRules(config *rawMainDefinition) []*validator.FieldRules {
	return []*validator.FieldRules{
		validator.Field(&config.Version, validator.Required, validator.By(versionValidator)),
		validator.Field(&config.Config, validator.By(configValidator)),
		validator.Field(&config.Variables, validator.By(variablesValidator)),
		validator.Field(&config.Branches,
			validator.Required,
			validator.Each(validator.Required, validator.By(branchValidator(config.Branches))),
		),
		validator.Field(&config.Tags,
			validator.Each(validator.By(tagValidator(config.Workflows, config.Branches))),
		),
		validator.Field(&config.Workflows,
			validator.Required,
			validator.Each(validator.Required, validator.By(workflowMapValidator)),
		),
	}
}
