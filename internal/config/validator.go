package config

import (
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

// ValidateConfig takes a raw configuration and validates it section by section
func validateWorkflowDefinition(workflowDefinition *rawWorkflowDefinition) error {
	if workflowDefinition.Flowit == nil {
		return errors.New("Workflow definition must contain a main 'flowit' key")
	}
	return validator.ValidateStruct(workflowDefinition.Flowit, mainDefinitionFieldRules(workflowDefinition.Flowit)...)
}

func mainDefinitionFieldRules(mainDefinition *rawMainDefinition) []*validator.FieldRules {
	return []*validator.FieldRules{
		validator.Field(&mainDefinition.Version, validator.Required, validator.By(versionValidator)),
		validator.Field(&mainDefinition.Config, validator.By(configValidator)),
		validator.Field(&mainDefinition.Variables, validator.By(variablesValidator)),
		validator.Field(&mainDefinition.Branches,
			validator.Required,
			validator.Each(validator.Required, validator.By(branchValidator(mainDefinition.Branches))),
		),
		validator.Field(&mainDefinition.Tags,
			validator.Each(validator.By(tagValidator(mainDefinition.Workflows, mainDefinition.Branches))),
		),
		validator.Field(&mainDefinition.Workflows,
			validator.Required,
			validator.Each(validator.Required, validator.By(workflowMapValidator)),
		),
	}
}
