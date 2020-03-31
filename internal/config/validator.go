package config

import (
	"fmt"
	"reflect"
	"runtime"

	validator "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

// ValidateConfig takes a viper configuration and validates it section by section
func validateConfig(flowit *rawFlowitConfig) error {

	var customValidators = map[string]interface{}{
		"flowitversion":           versionValidator,
		"flowitconfigshell":       shellValidator,
		"flowitbranches":          branchesValidator,
		"flowitbranchid":          branchIDValidator,
		"flowitbranchname":        branchNameValidator,
		"flowitbranchpreffix":     branchPreffixValidator,
		"flowitbranchsuffix":      branchSuffixValidator,
		"flowitbranchtransitions": transitionsValidator,
	}

	validator.SetFieldsRequiredByDefault(true)
	validator.SetNilPtrAllowedByRequired(false)

	if err := registerCustomValidators(customValidators); err != nil {
		return errors.Wrap(err, "Error registering validators")
	}

	// TODO: The error message when we return false is mostly unreadable. Can we change it?
	if _, err := validator.ValidateStruct(flowit); err != nil {
		return errors.Wrap(err, "Validation error")
	}
	return nil
}

func registerCustomValidators(validators map[string]interface{}) error {
	for k, v := range validators {
		switch v := v.(type) {
		case func(string) bool:
			validator.TagMap[k] = v
		case func(interface{}, interface{}) bool:
			validator.CustomTypeTagMap.Set(k, validator.CustomTypeValidator(v))
		default:
			return errors.New(
				fmt.Sprintf("Function %s is not a valid validator",
					runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()))
		}
	}
	return nil
}
