package configs

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func validateViperConfig(viperConfig *Viper) error {

	flowit, err := unmarshallConfig(viperConfig)
	if err != nil {
		return fmt.Errorf("Validation error: %w", err)
	}

	valid.SetFieldsRequiredByDefault(true)
	valid.SetNilPtrAllowedByRequired(false)

	registerCustomValidators()

	// TODO: The error message when we return false is mostly unreadable. Can we change it?
	_, err = valid.ValidateStruct(flowit)
	if err != nil {
		return fmt.Errorf("Validation error: %w", err)
	}
	return nil
}
