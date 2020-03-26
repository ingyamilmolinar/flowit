package configs

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func unmarshallConfig(viperConfig *Viper) (Flowit, error) {

	var flowit Flowit

	config := func(c *mapstructure.DecoderConfig) {
		c.ErrorUnused = true
		c.WeaklyTypedInput = false
	}

	err := viper.UnmarshalKey("flowit", &flowit, config)
	if err != nil {
		return flowit, fmt.Errorf("Unmarshalling error: %w", err)
	}
	return flowit, nil
}
