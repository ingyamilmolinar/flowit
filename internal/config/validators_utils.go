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
			append(
				commonNamingRules(),
				validator.NewStringRule(doesNotContainVariables, "Identifier contains variable references"),
			)...,
		)
	default:
		return errors.New("Invalid identifier type. Got " + reflect.TypeOf(identifier).Name())
	}
}

func validName(name interface{}) error {
	switch name := name.(type) {
	case *string, string:
		return validator.Validate(name,
			commonNamingRules()...,
		)
	default:
		return errors.New("Invalid name type. Got " + reflect.TypeOf(name).Name())
	}
}

func commonNamingRules() []validator.Rule {
	return []validator.Rule{
		is.PrintableASCII,
		validator.NewStringRule(doesNotContainWhiteSpace, "Field contains whitespaces"),
	}
}

func doesNotContainWhiteSpace(str string) bool {
	if ok, _ := regexp.Match(`\s`, []byte(str)); ok {
		return false
	}
	return true
}

func doesNotContainVariables(str string) bool {
	if ok, _ := regexp.Match(`\$<.*>`, []byte(str)); ok {
		return false
	}
	return true
}
