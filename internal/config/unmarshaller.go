package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func unmarshallConfig(v *viper.Viper) (*rawFlowitConfig, error) {

	var flowit rawFlowitConfig

	config := func(c *mapstructure.DecoderConfig) {
		c.ErrorUnused = true
		c.WeaklyTypedInput = false
	}

	if err := (*v).UnmarshalKey("flowit", &flowit, config); err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	return &flowit, nil
}
