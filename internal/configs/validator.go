package configs

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func validateViperConfig(flowit Flowit) error {

	valid.SetFieldsRequiredByDefault(true)
	valid.SetNilPtrAllowedByRequired(false)

	registerCustomValidators(map[string]interface{}{
		"flowitversion":           versionValidator,
		"flowitconfigshell":       shellValidator,
		"flowitbranches":          branchesValidator,
		"flowitbranchid":          branchIDValidator,
		"flowitbranchname":        branchNameValidator,
		"flowitbranchpreffix":     branchPreffixValidator,
		"flowitbranchsuffix":      branchSuffixValidator,
		"flowitbranchtransitions": transitionsValidator,
	})

	// TODO: The error message when we return false is mostly unreadable. Can we change it?
	_, err := valid.ValidateStruct(flowit)
	if err != nil {
		return fmt.Errorf("Validation error: %w", err)
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
			return fmt.Errorf("Function %s is not a valid validator", v)
		}
	}
	return nil
}
