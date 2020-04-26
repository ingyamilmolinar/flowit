package config

import (
	"reflect"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func stageMapValidator(workflowTypeMap interface{}) error {
	switch workflowTypeMap := workflowTypeMap.(type) {
	case rawWorkflowType:
		for workflowType, stages := range workflowTypeMap {
			if err := validator.Validate(workflowType,
				validator.Required,
				validator.By(workflowTypeValidator)); err != nil {
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
		return errors.New("Invalid workflow.stages type. Got " + reflect.TypeOf(workflowTypeMap).Name())
	}
}

// POST MVP: We may need to validate our variable syntax
func workflowTypeValidator(workflowType interface{}) error {
	switch workflowType := workflowType.(type) {
	case string:
		return nil
	default:
		return errors.New("Invalid workflow.stages workflow type. Got " + reflect.TypeOf(workflowType).Name())
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
		return errors.New("Invalid workflow.stages workflow type. Got " + reflect.TypeOf(stage).Name())
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
	case string:
		return nil
	default:
		return errors.New("Invalid workflow.stages args type. Got " + reflect.TypeOf(args).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageConditionsValidator(conditions interface{}) error {
	switch conditions := conditions.(type) {
	//TODO: Why not? case []string:
	case []interface{}:
		for _, condition := range conditions {
			_, ok := condition.(string)
			if !ok {
				return errors.New("Invalid workflow.stages condition type. Got " + reflect.TypeOf(condition).Name())
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow.stages conditions type. Got " + reflect.TypeOf(conditions).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageActionsValidator(actions interface{}) error {
	switch actions := actions.(type) {
	//TODO: Why not? case []string:
	case []interface{}:
		for _, action := range actions {
			_, ok := action.(string)
			if !ok {
				return errors.New("Invalid workflow.stages action type. Got " + reflect.TypeOf(action).Name())
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow.stages actions type. Got " + reflect.TypeOf(actions).Name())
	}
}
