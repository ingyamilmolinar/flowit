package config

import (
	validator "github.com/asaskevich/govalidator"
)

func branchesValidator(current interface{}, parent interface{}) bool {
	switch current := current.(type) {
	case []*rawBranch:
		branches := current
		if len(branches) == 0 {
			return false
		}
		for _, branch := range branches {
			// TODO: This is hiding internal errors
			if ok, _ := validator.ValidateStruct(*branch); !ok {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func branchIDValidator(id string) bool {
	return validIdentifier(id)
}

func branchNameValidator(name string) bool {
	return validIdentifier(name)
}

func branchPreffixValidator(preffix string) bool {
	return validIdentifier(preffix)
}

func branchSuffixValidator(suffix string) bool {
	return validIdentifier(suffix)
}

func validIdentifier(str string) bool {
	return str != "" && !validator.HasWhitespace(str) && validator.IsPrintableASCII(str)
}

func transitionsValidator(current interface{}, parent interface{}) bool {
	isParentBranchEternal := false
	switch parent := parent.(type) {
	case rawBranch:
		branch := parent
		if branch.Eternal == nil {
			return false
		}
		isParentBranchEternal = *branch.Eternal
	default:
		return false
	}

	switch current := current.(type) {
	case []*rawTransition:
		transitions := current
		if (!isParentBranchEternal && len(transitions) == 0) || (isParentBranchEternal && len(transitions) > 0) {
			return false
		}
	default:
		return false
	}
	return true
}
