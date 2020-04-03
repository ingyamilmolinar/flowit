package config

import (
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

// ValidateConfig takes a viper configuration and validates it section by section
func validateConfig(flowit *rawFlowitConfig) error {
	return validator.ValidateStruct(flowit, configFieldRules(flowit)...)
}

func configFieldRules(flowit *rawFlowitConfig) []*validator.FieldRules {
	return []*validator.FieldRules{
		validator.Field(&flowit.Version, validator.Required, validator.By(versionValidator)),
		validator.Field(&flowit.Config, validator.By(configValidator)),
		validator.Field(&flowit.Workflow, workflowRules()...),
	}
}
