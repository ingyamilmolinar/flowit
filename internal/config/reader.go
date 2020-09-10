package config

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func readWorkflowDefinition(fileLocation string) (*viper.Viper, error) {

	fileDir := filepath.Dir(fileLocation)
	fileExt := filepath.Ext(fileLocation)
	filenameWithoutExt := strings.ReplaceAll(filepath.Base(fileLocation), filepath.Ext(fileLocation), "")

	viper := viper.New()
	viper.SetConfigType(strings.ReplaceAll(fileExt, ".", ""))
	viper.SetConfigName(filenameWithoutExt)
	viper.AddConfigPath(fileDir)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Workflow definition read error")
	}
	return viper, nil
}
