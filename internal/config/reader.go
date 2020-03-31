package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configType = "yaml"

func readConfig(configName string, configLocation string) (*viper.Viper, error) {
	viper := viper.New()
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Config reading error")
	}
	return viper, nil
}
