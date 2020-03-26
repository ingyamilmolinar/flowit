package configs

import (
	"fmt"

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

// ReadConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ParseConfig(configName string, configLocation string) (*Viper, error) {

	// TODO: Hash parsed and validated config and verify if it changed or not

	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Fatal error reading config file: %w", err)
	}
	viper := Viper{viper.GetViper()}
	err = validateViperConfig(&viper)
	if err != nil {
		return nil, err
	}
	return &viper, nil
}
