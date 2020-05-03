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

func stageValidator(stage interface{}) error {
	switch stage := stage.(type) {
	case rawStage:
		var foundID, foundActions bool
		for k, v := range stage {
			switch k {
			case "id":
				foundID = true
				if err := validator.Validate(v, validator.By(stageIDValidator)); err != nil {
					return err
				}
			case "args":
				if err := validator.Validate(v, validator.By(stageArgsValidator)); err != nil {
					return err
				}
			case "conditions":
				if err := validator.Validate(v, validator.By(stageConditionsValidator)); err != nil {
					return err
				}
			case "actions":
				foundActions = true
				if err := validator.Validate(v, validator.By(stageActionsValidator)); err != nil {
					return err
				}
			default:
				return errors.New("Invalid workflow stage: " + k + " section is not valid")
			}
		}
		if !foundID {
			return errors.New("Invalid workflow stage: non optional section 'id' was not found")
		}
		if !foundActions {
			return errors.New("Invalid workflow stage: non optional section 'actions' was not found")
		}
	default:
		return errors.New("Invalid workflow stages type. Got " + reflect.TypeOf(stage).Name())
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
	// TODO: Abstract repeated code
	case []interface{}:
		for _, action := range args {
			_, ok := action.(string)
			if !ok {
				return errors.New("Invalid workflow stages arg type. Got " + reflect.TypeOf(action).Name())
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow stages args type. Got " + reflect.TypeOf(args).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageConditionsValidator(conditions interface{}) error {
	switch conditions := conditions.(type) {
	// TODO: Abstract repeated code
	case []interface{}:
		for _, condition := range conditions {
			_, ok := condition.(string)
			if !ok {
				return errors.New("Invalid workflow stages condition type. Got " + reflect.TypeOf(condition).Name())
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow stages conditions type. Got " + reflect.TypeOf(conditions).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageActionsValidator(actions interface{}) error {
	switch actions := actions.(type) {
	// TODO: Abstract repeated code
	case []interface{}:
		for _, action := range actions {
			_, ok := action.(string)
			if !ok {
				return errors.New("Invalid workflow stages action type. Got " + reflect.TypeOf(action).Name())
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow stages actions type. Got " + reflect.TypeOf(actions).Name())
	}
}