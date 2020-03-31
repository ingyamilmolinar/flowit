package utils

import (
	"fmt"
	"os"
)

// ExitIfErr logs and panics if an error exists
func ExitIfErr(err error) {
	if err != nil {
		const exitStatus = 1
		GetLogger().Error(fmt.Sprintf("%+v", err))
		os.Exit(exitStatus)
	}
}
