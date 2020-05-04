package config

import (
	"reflect"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func workflowMapValidator(workflowMap interface{}) error {
	switch workflowMap := workflowMap.(type) {
	case rawWorkflow:
		for workflow, stages := range workflowMap {
			if err := validator.Validate(workflow,
				validator.Required,
				validator.By(workflowValidator)); err != nil {
				return err
			}
			if err := validator.Validate(stages,
				validator.Required,
				validator.By(stagesValidator),
				validator.Each(
					validator.Required,
					validator.By(stageValidator))); err != nil {
				return err
			}

		}
		return nil
	default:
		return errors.New("Invalid workflows type. Got " + reflect.TypeOf(workflowMap).Name())
	}
}

// POST MVP: We may need to validate our variable syntax
func workflowValidator(workflow interface{}) error {
	switch workflow := workflow.(type) {
	case string:
		return nil
	default:
		return errors.New("Invalid workflows workflow type. Got " + reflect.TypeOf(workflow).Name())
	}
}

func stagesValidator(stages interface{}) error {
	switch stages := stages.(type) {
	case []*rawStage:
		var foundStart, foundFinish bool
		for _, stage := range stages {
			if *stage.ID == "start" {
				foundStart = true
			}
			if *stage.ID == "finish" {
				foundFinish = true
			}
		}
		if !foundStart || !foundFinish {
			return errors.New("Invalid workflow stages: 'start' and 'finish' stages are required")
		}
	default:
		return errors.New("Invalid workflow stages type. Got " + reflect.TypeOf(stages).Name())
	}
	return nil
}

func stageValidator(stage interface{}) error {
	switch stage := stage.(type) {
	case rawStage:
		if err := validator.Validate(stage.ID, validator.Required, validator.By(stageIDValidator)); err != nil {
			return err
		}
		if err := validator.Validate(stage.Args, validator.By(stageArgsValidator)); err != nil {
			return err
		}
		if err := validator.Validate(stage.Conditions, validator.By(stageConditionsValidator)); err != nil {
			return err
		}
		if err := validator.Validate(stage.Actions, validator.Required, validator.By(stageActionsValidator)); err != nil {
			return err
		}
	default:
		return errors.New("Invalid workflow stage type. Got " + reflect.TypeOf(stage).Name())
	}
	return nil
}

func stageIDValidator(id interface{}) error {
	return validator.Validate(id,
		validator.Required,
		validator.By(validIdentifier))
}

// TODO: We may need to validate our variable syntax
func stageArgsValidator(args interface{}) error {
	switch args := args.(type) {
	case []*string:
		return nil
	default:
		return errors.New("Invalid workflow stages args type. Got " + reflect.TypeOf(args).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageConditionsValidator(conditions interface{}) error {
	switch conditions := conditions.(type) {
	case []*string:
		return nil
	default:
		return errors.New("Invalid workflow stages conditions type. Got " + reflect.TypeOf(conditions).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageActionsValidator(actions interface{}) error {
	switch actions := actions.(type) {
	case []*string:
		return nil
	default:
		return errors.New("Invalid workflow stages actions type. Got " + reflect.TypeOf(actions).Name())
	}
}
