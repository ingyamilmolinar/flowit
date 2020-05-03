package config

import (
	"reflect"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func variablesValidator(variables interface{}) error {
	switch variables := variables.(type) {
	case *rawVariables:
		// variables section is not mandatory
		if variables == nil {
			return nil
		}
		if len(*variables) == 0 {
			return errors.New("Variables can not be both present on the configuration AND empty")
		}
		for _, variableValue := range *variables {
			if err := validator.Validate(variableValue, validator.By(variableValueValidator)); err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("Invalid variables type. Got " + reflect.TypeOf(variables).Name())
	}
}

// TODO: We may need to validate user defined variables
func variableValueValidator(variableValue interface{}) error {
	return nil
}
