package config

import (
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

// ProcessFlowitConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ProcessFlowitConfig(configName string, configLocation string) (*FlowitConfig, error) {

	// TODO: Hash parsed and validated config and verify if it changed or not?
	viper, err := readConfig(configName, configLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Config reading error")
	}

	rawConfig, err := unmarshallConfig(viper)
	if err != nil {
		return nil, errors.Wrap(err, "Config unmarshalling error")
	}

	if err := validateConfig(rawConfig); err != nil {
		return nil, errors.Wrap(err, "Config validation error")
	}

	//TODO: Include defaults on rawConfig, use viper
	var config FlowitConfig
	if err := utils.DeepCopy(rawConfig, &config); err != nil {
		return nil, errors.Wrap(err, "Config copying error")
	}
	return &config, nil
}
