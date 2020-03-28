package configs

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ValidateViperConfig takes a viper configuration and validates it section by section
func unmarshallConfig(v *viper.Viper) (*Flowit, error) {

	var flowit Flowit

	config := func(c *mapstructure.DecoderConfig) {
		c.ErrorUnused = true
		c.WeaklyTypedInput = false
	}

	err := (*v).UnmarshalKey("flowit", &flowit, config)
	if err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	return &flowit, nil
}
