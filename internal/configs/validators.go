package configs

import (
	valid "github.com/asaskevich/govalidator"
)

var supportedFlowitVersions = []string{"0.1"}

func registerCustomValidators() {
	valid.TagMap["flowitversion"] = validateVersion
	valid.CustomTypeTagMap.Set("flowitbranches", valid.CustomTypeValidator(validateBranches))
}

func validateVersion(str string) bool {
	for _, supportedVersion := range supportedFlowitVersions {
		if supportedVersion == str {
			return true
		}
	}
	return false
}
