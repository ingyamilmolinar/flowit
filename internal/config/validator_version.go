package config

import (
	"reflect"

	"github.com/pkg/errors"
)

func versionValidator(version interface{}) error {
	switch version := version.(type) {
	case *string:
		var supportedVersions = []string{"0.1"}
		found := false
		for _, supportedVersion := range supportedVersions {
			if *version == supportedVersion {
				found = true
			}
		}
		if !found {
			return errors.New("Unsupported workflow definition version")
		}
	default:
		return errors.New("Invalid version type. Got " + reflect.TypeOf(version).Name())
	}
	return nil
}
