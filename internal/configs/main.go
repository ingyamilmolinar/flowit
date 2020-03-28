package configs

import (
	"github.com/pkg/errors"
)

// ProcessFlowitConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ProcessFlowitConfig(configName string, configLocation string) (*Flowit, error) {

	// TODO: Hash parsed and validated config and verify if it changed or not?
	viper, err := readConfig(configName, configLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Config reading error")
	}

	flowit, err := unmarshallConfig(viper)
	if err != nil {
		return nil, errors.Wrap(err, "Config unmarshalling error")
	}

	err = validateConfig(flowit)
	if err != nil {
		return flowit, errors.Wrap(err, "Config validation error")
	}
	return flowit, nil
}
