package configs

import (
	"github.com/spf13/viper"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var logger = utils.GetLogger()

// ValidateConfig takes a viper configuration and validates it section by section
func ValidateConfig(viperConfig *viper.Viper) {
	version := viperConfig.GetString("flowit.version")
	config := viperConfig.Get("flowit.config").([]interface{})
	variables := viperConfig.Get("flowit.variables").([]interface{})
	workflow := viperConfig.Get("flowit.workflow").(map[string]interface{})
	logger.Info(version)
	logger.Info(config[0].(map[interface{}]interface{})["abort-on-failed-action"])
	logger.Info(variables[0].(map[interface{}]interface{})["circleci-username"])
	logger.Info(workflow["branches"].([]interface{})[0].(map[interface{}]interface{}))
}
