package config

import (
	"reflect"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func workflowValidator(stateMachines []*rawStateMachine) func(interface{}) error {
	return func(workflow interface{}) error {
		switch workflow := workflow.(type) {
		case rawWorkflow:
			if workflow.ID == nil {
				return errors.New("Workflow ID is nil")
			}
			if err := validator.Validate(*workflow.ID,
				validator.Required,
				validator.By(workflowIDValidator)); err != nil {
				return errors.WithStack(err)
			}
			if workflow.StateMachine == nil {
				return errors.New("Workflow StateMachine is nil")
			}
			if err := validator.Validate(*workflow.StateMachine,
				validator.Required,
				validator.NewStringRule(
					workflowStateMachineIDValidator(stateMachines),
					"Workflow State Machine ID is not a valid state machine")); err != nil {
				return errors.WithStack(err)
			}
			if err := validator.Validate(workflow.Stages,
				validator.Required,
				validator.By(workflowStagesValidator(*workflow.StateMachine, stateMachines)),
				validator.Each(
					validator.Required,
					validator.By(workflowStageValidator))); err != nil {
				return errors.WithStack(err)
			}
			return nil
		default:
			return errors.New("Invalid workflow type. Got " + reflect.TypeOf(workflow).Name())
		}
	}
}

func workflowIDValidator(workflowID interface{}) error {
	return validIdentifier(workflowID)
}

func workflowStateMachineIDValidator(stateMachines []*rawStateMachine) func(string) bool {
	return func(stateMachineID string) bool {
		found := false
		for _, stateMachine := range stateMachines {
			if stateMachine.ID == nil {
				continue
			}
			if *stateMachine.ID == stateMachineID {
				found = true
				break
			}
		}
		return found
	}
}

func workflowStagesValidator(workflowStateMachineID string, stateMachines []*rawStateMachine) func(interface{}) error {
	return func(stages interface{}) error {
		switch stages := stages.(type) {
		case []*rawStage:
			var workflowStateMachine *rawStateMachine
			for _, stateMachine := range stateMachines {
				if stateMachine.ID == nil {
					continue
				}
				if workflowStateMachineID == *stateMachine.ID {
					workflowStateMachine = stateMachine
				}
			}
			if workflowStateMachine == nil {
				return errors.New("Invalid state machine ID: " + workflowStateMachineID)
			}

			foundStageCounter := 0
			for _, stage := range stages {
				foundStageID := false
				for _, stateMachineStages := range workflowStateMachine.Stages {
					if stage.ID == nil || stateMachineStages == nil {
						continue
					}
					if *stage.ID == *stateMachineStages {
						foundStageID = true
						foundStageCounter++
						break
					}
				}
				if !foundStageID && stage.ID != nil && workflowStateMachine.ID != nil {
					return errors.New("Stage with ID: " + *stage.ID +
						" is not a valid " + *workflowStateMachine.ID + " state machine stage")
				}
			}
			if foundStageCounter != len(workflowStateMachine.Stages) {
				return errors.New("Some " + *workflowStateMachine.ID +
					" state machine stages are missing in workflow")
			}
		default:
			return errors.New("Invalid workflow stages type. Got " + reflect.TypeOf(stages).Name())
		}
		return nil
	}
}

func workflowStageValidator(stage interface{}) error {
	switch stage := stage.(type) {
	case rawStage:
		if err := validator.Validate(stage.ID, validator.Required, validator.By(stageIDValidator)); err != nil {
			return errors.WithStack(err)
		}
		if err := validator.Validate(stage.Args, validator.By(stageArgsValidator)); err != nil {
			return errors.WithStack(err)
		}
		if err := validator.Validate(stage.Conditions, validator.By(stageConditionsValidator)); err != nil {
			return errors.WithStack(err)
		}
		if err := validator.Validate(stage.Actions, validator.Required, validator.By(stageActionsValidator)); err != nil {
			return errors.WithStack(err)
		}
	default:
		return errors.New("Invalid workflow stage type. Got " + reflect.TypeOf(stage).Name())
	}
	return nil
}

func stageIDValidator(id interface{}) error {
	return validIdentifier(id)
}

func stageArgsValidator(args interface{}) error {
	switch args := args.(type) {
	case []*string:
		for _, arg := range args {
			if !utils.IsValidVariableDeclaration(*arg) {
				return errors.New("Invalid workflow stage argument: " + (*arg))
			}
		}
		return nil
	default:
		return errors.New("Invalid workflow stage arguments type. Got " + reflect.TypeOf(args).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageConditionsValidator(conditions interface{}) error {
	switch conditions := conditions.(type) {
	case []*string:
		return nil
	default:
		return errors.New("Invalid workflow stage conditions type. Got " + reflect.TypeOf(conditions).Name())
	}
}

// TODO: We may need to validate our variable syntax
func stageActionsValidator(actions interface{}) error {
	switch actions := actions.(type) {
	case []*string:
		return nil
	default:
		return errors.New("Invalid workflow stage actions type. Got " + reflect.TypeOf(actions).Name())
	}
}
