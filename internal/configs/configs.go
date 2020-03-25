package configs

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/yamil-rivera/flowit/internal/utils"
)

// ReadConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ReadConfig(configName string, configLocation string) *viper.Viper {
	var logger = utils.GetLogger()

	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Errorf("Fatal error config file: %s", err)
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	return viper.GetViper()
}
