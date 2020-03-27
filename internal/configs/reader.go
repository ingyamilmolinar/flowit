package configs

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configType = "yaml"

func readConfig(configName string, configLocation string) error {
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Config reading error")
	}
	return nil
}
