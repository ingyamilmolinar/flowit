package config

import (
	"os/exec"
	"reflect"
	"strings"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func configValidator(config interface{}) error {
	switch config := config.(type) {
	case *rawConfig:
		// TODO: Fail on empty map not on missing map. Apply fix to all optional sections
		// config is optional
		if config == nil {
			return nil
		}
		return validator.Validate(config.Shell, validator.By(shellValidator))
	default:
		return errors.New("Invalid config type. Got " + reflect.TypeOf(config).Name())
	}
}

// TODO: Validate command and args before executing
func shellValidator(shell interface{}) error {
	switch shell := shell.(type) {
	case *string:
		if shell != nil {
			cmds := strings.Split(*shell, " ")
			/* #gosec */
			cmd := exec.Command(cmds[0], cmds[1:]...)
			_, err := cmd.Output()
			return err
		}
	default:
		return errors.New("Invalid config shell type. Got " + reflect.TypeOf(shell).Name())
	}
	return nil
}
