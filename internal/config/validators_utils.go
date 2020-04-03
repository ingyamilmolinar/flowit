package config

import (
	"reflect"
	"regexp"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pkg/errors"
)

// TODO: This needs more thought
func validIdentifier(identifier interface{}) error {
	switch identifier := identifier.(type) {
	case *string, string:
		return validator.Validate(identifier,
			is.PrintableASCII,
			validator.NewStringRule(doesNotContainWhiteSpace, "Identifier contains whitespaces"),
		)
	default:
		return errors.New("Invalid identifier type. Got " + reflect.TypeOf(identifier).Name())
	}
}

func doesNotContainWhiteSpace(str string) bool {
	if ok, _ := regexp.Match("\\s", []byte(str)); ok {
		return false
	}
	return true
}
