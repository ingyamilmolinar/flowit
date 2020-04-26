package config

import (
	"reflect"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func workflowValidator(workflow interface{}) error {
	switch workflow := workflow.(type) {
	case *rawWorkflow:
		return validator.ValidateStruct(workflow,
			validator.Field(&workflow.Branches,
				validator.Required,
				validator.Each(validator.Required, validator.By(branchValidator(workflow.Branches))),
			),
			validator.Field(&workflow.Tags,
				validator.Each(validator.By(tagValidator(workflow.Stages, workflow.Branches))),
			),
			validator.Field(&workflow.Stages,
				validator.Required,
				validator.Each(validator.Required, validator.By(stageMapValidator)),
			),
		)
	default:
		return errors.New("Invalid flowit.workflow type. Got " + reflect.TypeOf(workflow).Name())
	}
}
