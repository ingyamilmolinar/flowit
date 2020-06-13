package config

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func versionValidator(version interface{}) error {
	switch version := version.(type) {
	case *string:
		var supportedVersions = []string{"0.1"}
		if found := utils.FindStringInArray(*version, supportedVersions); !found {
			return errors.New("Unsupported workflow definition version")
		}
	default:
		return errors.New("Invalid version type. Got " + reflect.TypeOf(version).Name())
	}
	return nil
}
