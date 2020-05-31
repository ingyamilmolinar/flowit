package utils

import "errors"

// OptionalString is the data type that wraps a string in an optional
type OptionalString struct {
	str   string
	isSet bool
}

// NewString receives a string and returns a filled optional ready to be unwrapped
func NewString(str string) OptionalString {
	return OptionalString{
		str:   str,
		isSet: true,
	}
}

// Get works on a pointer to OptionalString and returns the string wrapped in the optional or
// returns an error in case the optional is empty
func (optional *OptionalString) Get() (string, error) {
	if !optional.isSet {
		return optional.str, errors.New("optional value is not set")
	}
	return optional.str, nil
}
