package configs

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
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Config reading error")
	}
	return viper, nil
}
