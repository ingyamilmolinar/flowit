package config

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	validator "github.com/go-ozzo/ozzo-validation/v4"
)

func branchValidator(branches []*rawBranch) validator.RuleFunc {
	return func(branch interface{}) error {
		switch branch := branch.(type) {
		case rawBranch:
			return validator.ValidateStruct(&branch,
				validator.Field(&branch.ID, validator.By(branchIDValidator)),
				validator.Field(&branch.Name, validator.By(branchNameValidator)),
				validator.Field(&branch.Prefix, validator.By(branchPreffixValidator)),
				validator.Field(&branch.Suffix, validator.By(branchSuffixValidator)),
				validator.Field(&branch.Eternal, validator.NotNil),
				validator.Field(&branch.Protected, validator.NotNil),
				validator.Field(&branch.Transitions,
					validator.By(branchTransitionsValidator(branch.Eternal)),
					validator.Each(validator.By(branchTransitionValidator(branches)))),
			)
		default:
			return errors.New("Invalid workflow.branch type. Got " + reflect.TypeOf(branch).Name())
		}
	}
}

func branchIDValidator(id interface{}) error {
	return validator.Validate(id,
		validator.Required,
		validator.By(validIdentifier),
	)
}

// TODO: We may need to validate our variable syntax
func branchNameValidator(name interface{}) error {
	return validator.Validate(name,
		validator.Required,
		validator.By(validName),
	)
}

// TODO: We may need to validate our variable syntax
func branchPreffixValidator(preffix interface{}) error {
	return validName(preffix)
}

// TODO: We may need to validate our variable syntax
func branchSuffixValidator(suffix interface{}) error {
	return validName(suffix)
}

func branchTransitionsValidator(eternal *bool) validator.RuleFunc {
	return func(transitions interface{}) error {
		switch transitions := transitions.(type) {
		case []*rawTransition:
			if !*eternal && len(transitions) == 0 {
				return errors.New("Invalid branch transitions: Transitions must be specified for non eternal branches")
			} else if *eternal && len(transitions) > 0 {
				return errors.New("Invalid branch transitions: Transitions must not be specified for eternal branches")
			}
			// TODO: Check for repeated transitions
			return nil
		default:
			return errors.New("Invalid branch.transitions type. Got " + reflect.TypeOf(transitions).Name())
		}
	}
}

func branchTransitionValidator(branches []*rawBranch) validator.RuleFunc {
	return func(transition interface{}) error {
		switch transition := transition.(type) {
		case rawTransition:
			return validator.ValidateStruct(&transition,
				validator.Field(&transition.From,
					validator.Required,
					validator.By(branchTransitionFromValidator(branches))),
				validator.Field(&transition.To,
					validator.Required,
					validator.Each(
						validator.Required,
						validator.By(branchTransitionToValidator(branches)))),
			)
		default:
			return errors.New("Invalid branch.transition type. Got " + reflect.TypeOf(transition).Name())
		}
	}
}

func branchTransitionFromValidator(branches []*rawBranch) validator.RuleFunc {
	return func(from interface{}) error {
		switch from := from.(type) {
		case *string:
			found := isBranchDefined(*from, branches)
			if !found {
				return errors.New("Invalid branch transition: " + *from + " is not a defined branch")
			}
		default:
			return errors.New("Invalid branch.transition.from type. Got " + reflect.TypeOf(from).Name())
		}
		return nil
	}
}

func branchTransitionToValidator(branches []*rawBranch) validator.RuleFunc {
	return func(to interface{}) error {
		switch to := to.(type) {
		case string:
			split := strings.Split(to, ":")
			if len(split) != 2 {
				return errors.New("Invalid branch transition: 'to' should be of the form '<branch>:<local|remote>'")
			}
			branch := split[0]
			if !isBranchDefined(branch, branches) {
				return errors.New("Invalid branch transition: " + branch + " is not a defined branch")
			}
			option := split[1]
			if option != "local" && option != "remote" {
				return errors.New("Invalid branch transition: " + branch + " option should be 'local' or 'remote'")
			}
		default:
			return errors.New("Invalid branch.transition.to type. Got " + reflect.TypeOf(to).Name())
		}
		return nil
	}
}

func isBranchDefined(branch string, branches []*rawBranch) bool {
	for _, definedBranch := range branches {
		if definedBranch != nil && definedBranch.ID != nil && branch == *definedBranch.ID {
			return true
		}
	}
	return false
}
