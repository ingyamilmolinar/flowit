package config

import (
	"reflect"

	"github.com/pkg/errors"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func tagValidator(workflows []*rawWorkflowType, branches []*rawBranch) validator.RuleFunc {
	return func(tag interface{}) error {
		switch tag := tag.(type) {
		case rawTag:
			return validator.ValidateStruct(&tag,
				validator.Field(&tag.ID, validator.By(tagIDValidator)),
				validator.Field(&tag.Format, validator.By(tagFormatValidator)),
				validator.Field(&tag.Stages, validator.By(tagStagesValidator(workflows))),
				validator.Field(&tag.Branches, validator.Each(validator.By(tagBranchesValidator(branches)))),
			)
		default:
			return errors.New("Invalid workflow.tag type. Got " + reflect.TypeOf(tag).Name())
		}
	}
}

func tagIDValidator(id interface{}) error {
	return validator.Validate(id,
		validator.Required,
		validator.By(validIdentifier),
	)
}

// TODO: More thought needs to be put into this
func tagFormatValidator(format interface{}) error {
	return validator.Validate(format,
		is.PrintableASCII,
	)
}

func tagStagesValidator(workflows []*rawWorkflowType) validator.RuleFunc {
	return func(stageMap interface{}) error {
		switch stageMap := stageMap.(type) {
		case *map[string][]*string:
			if stageMap == nil {
				return nil
			}
			for workflow, stages := range *stageMap {
				foundWorkflow := false
				for _, workflowStages := range workflows {
					if workflowStages == nil {
						return nil
					}
					// TODO: Does not work if workflow has variables
					if _, ok := (*workflowStages)[workflow]; !ok {
						continue
					}
					foundWorkflow = true
					foundStages := areStagesDefined(stages, (*workflowStages)[workflow])
					if !foundStages {
						return errors.New("Invalid tag stages: Stage under workflow " + workflow + " is not defined")
					}
				}
				if !foundWorkflow {
					return errors.New("Invalid tag workflow: workflow " + workflow + " is not defined")
				}
			}
		default:
			return errors.New("Invalid tag.stages type. Got " + reflect.TypeOf(stageMap).Name())
		}
		return nil
	}
}

func areStagesDefined(stages []*string, definedStages []*rawStage) bool {
	foundStage := 0
	for _, definedStage := range definedStages {
		for _, stage := range stages {
			if stage == nil {
				return false
			}
			if stageID, ok := (*definedStage)["id"]; ok {
				if *stage == stageID {
					foundStage++
				}
			}
		}
	}
	return foundStage == len(stages)
}

func tagBranchesValidator(branches []*rawBranch) validator.RuleFunc {
	return func(branch interface{}) error {
		switch branch := branch.(type) {
		case string:
			if ok := isBranchDefined(branch, branches); !ok {
				return errors.New("Invalid tag branches: Branch " + branch + " is not defined")
			}
		default:
			return errors.New("Invalid tag.branches type. Got " + reflect.TypeOf(branch).Name())
		}
		return nil
	}
}