package configs

import (
	"fmt"

	"github.com/spf13/viper"
	logging "github.com/yamil-rivera/flowit/internal/utils"
)

// ReadConfig reads, parses the specified yaml configuration file and returns a map with the key/values
func ReadConfig(configName string) map[string]interface{} {
	var logger = logging.GetLogger()

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./samples/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Errorf("Fatal error config file: %s", err)
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	return viper.AllSettings()
}
