package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const configType = "yaml"

func readWorkflowDefinition(fileName string, fileLocation string) (*viper.Viper, error) {
	viper := viper.New()
	viper.SetConfigType(configType)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(fileLocation)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Workflow definition read error")
	}
	return viper, nil
}
