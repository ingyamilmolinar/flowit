package configs

import (
	"github.com/pkg/errors"
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

	err = validateConfig(rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Config validation error")
	}

	//TODO: Include defaults on rawConfig, use viper

	var config FlowitConfig
	err = deepCopy(rawConfig, &config)
	if err != nil {
		return nil, errors.Wrap(err, "Config copying error")
	}
	return &config, nil
}
