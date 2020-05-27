package config

import (
	"reflect"

	"github.com/pkg/errors"

	validator "github.com/go-ozzo/ozzo-validation/v4"
)

func tagValidator(workflows []*rawWorkflow, branches []*rawBranch) validator.RuleFunc {
	return func(tag interface{}) error {
		switch tag := tag.(type) {
		case rawTag:
			return validator.ValidateStruct(&tag,
				validator.Field(&tag.ID,
					validator.Required,
					validator.By(tagIDValidator)),
				validator.Field(&tag.Format,
					validator.Required,
					validator.By(tagFormatValidator)),
				validator.Field(&tag.Stages,
					validator.By(tagStagesValidator(workflows))),
				validator.Field(&tag.Branches,
					validator.Each(
						validator.NewStringRule(tagBranchesValidator(branches), "Invalid tag branches: Branch is not defined"))),
			)
		default:
			return errors.New("Invalid tag type. Got " + reflect.TypeOf(tag).Name())
		}
	}
}

func tagIDValidator(id interface{}) error {
	return validator.Validate(id,
		validator.By(validIdentifier),
	)
}

// TODO: More thought needs to be put into this
func tagFormatValidator(format interface{}) error {
	return validator.Validate(format,
		commonNamingRules()...,
	)
}

func tagStagesValidator(workflows []*rawWorkflow) validator.RuleFunc {
	return func(stageMap interface{}) error {
		switch stageMap := stageMap.(type) {
		case *rawStages:
			// Tag stages can be optional
			if stageMap == nil {
				return nil
			}
			for workflow, stages := range *stageMap {
				foundWorkflow := false
				for _, workflowStages := range workflows {
					if workflowStages == nil || workflowStages.ID == nil {
						continue
					}
					if *workflowStages.ID != workflow {
						continue
					}
					foundWorkflow = true
					foundStages := areStagesDefined(stages, workflowStages.Stages)
					if !foundStages {
						return errors.New("Invalid tag stages: Stage under workflow " + workflow + " is not defined")
					}
				}
				if !foundWorkflow {
					return errors.New("Invalid tag workflow: workflow " + workflow + " is not defined")
				}
			}
		default:
			return errors.New("Invalid tag stages type. Got " + reflect.TypeOf(stageMap).Name())
		}
		return nil
	}
}

func areStagesDefined(stages []*string, definedStages []*rawStage) bool {
	foundStage := 0
	for _, definedStage := range definedStages {
		for _, stage := range stages {
			if stage == nil || definedStage.ID == nil {
				return false
			}
			stageID := (*definedStage.ID)
			if *stage == stageID {
				foundStage++
			}
		}
	}
	return foundStage == len(stages)
}

func tagBranchesValidator(branches []*rawBranch) func(string) bool {
	return func(branch string) bool {
		return isBranchDefined(branch, branches)
	}
}
