package configs

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configType = "yaml"

// Viper is a custom type that embeds viper config. This was done to enhance viper with additional methods
type Viper struct {
	*viper.Viper
}

// GetSlice is a convenient way to get a slice of anything from a viper config. Idk why is not part of the viper API
// TODO: Should we open an issue/PR?
func (v *Viper) GetSlice(key string) []interface{} {
	return v.Get(key).([]interface{})
}

// ProcessFlowitConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ProcessFlowitConfig(configName string, configLocation string) (*Viper, error) {

	// TODO: Hash parsed and validated config and verify if it changed or not

	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Config reading error")
	}
	viper := Viper{viper.GetViper()}

	flowit, err := unmarshallConfig(&viper)
	if err != nil {
		return nil, errors.Wrap(err, "Config unmarshalling error")
	}

	err = validateViperConfig(flowit)
	if err != nil {
		return nil, errors.Wrap(err, "Config validation error")
	}
	return &viper, nil
}
