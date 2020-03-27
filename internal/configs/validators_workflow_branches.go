package configs

import (
	valid "github.com/asaskevich/govalidator"
)

func branchesValidator(current interface{}, parent interface{}) bool {
	switch current.(type) {
	case []*branch:
		branches := current.([]*branch)
		if len(branches) == 0 {
			return false
		}
		for _, branch := range branches {
			// TODO: This is hiding internal errors
			ok, _ := valid.ValidateStruct(*branch)
			if !ok {
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
	return str != "" && !valid.HasWhitespace(str) && valid.IsPrintableASCII(str)
}

func transitionsValidator(current interface{}, parent interface{}) bool {
	isParentBranchEternal := false
	switch parent.(type) {
	case branch:
		if parent.(branch).Eternal == nil {
			return false
		}
		isParentBranchEternal = *(parent.(branch).Eternal)
	default:
		return false
	}

	switch current.(type) {
	case []*transition:
		transitions := current.([]*transition)
		if (!isParentBranchEternal && len(transitions) == 0) || (isParentBranchEternal && len(transitions) > 0) {
			return false
		}
	default:
		return false
	}
	return true
}
