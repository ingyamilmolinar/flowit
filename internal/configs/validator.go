package configs

import (
	"fmt"
	"reflect"
	"runtime"

	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func validateViperConfig(flowit Flowit) error {

	valid.SetFieldsRequiredByDefault(true)
	valid.SetNilPtrAllowedByRequired(false)

	err := registerCustomValidators(map[string]interface{}{
		"flowitversion":           versionValidator,
		"flowitconfigshell":       shellValidator,
		"flowitbranches":          branchesValidator,
		"flowitbranchid":          branchIDValidator,
		"flowitbranchname":        branchNameValidator,
		"flowitbranchpreffix":     branchPreffixValidator,
		"flowitbranchsuffix":      branchSuffixValidator,
		"flowitbranchtransitions": transitionsValidator,
	})
	if err != nil {
		return errors.Wrap(err, "Error registering validators")
	}

	// TODO: The error message when we return false is mostly unreadable. Can we change it?
	_, err = valid.ValidateStruct(flowit)
	if err != nil {
		return errors.Wrap(err, "Validation error")
	}
	return nil
}

func registerCustomValidators(validators map[string]interface{}) error {
	for k, v := range validators {
		switch v.(type) {
		case func(string) bool:
			valid.TagMap[k] = v.(func(string) bool)
		case func(interface{}, interface{}) bool:
			valid.CustomTypeTagMap.Set(k, valid.CustomTypeValidator(v.(func(interface{}, interface{}) bool)))
		default:
			return errors.New(fmt.Sprintf("Function %s is not a valid validator", runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()))
		}
	}
	return nil
}
