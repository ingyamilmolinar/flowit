package utils

import (
	"fmt"
	"os"
)

// ExitIfErr logs and panics if an error exists
func ExitIfErr(err error) {
	if err != nil {
		GetLogger().Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
