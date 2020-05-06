package utils

import (
	"fmt"
	"os"
)

// OptionalExit logs and panics if an error exists
func OptionalExit(err error) {
	if err != nil {
		GetLogger().Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
