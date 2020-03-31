package utils

import (
	"path/filepath"
	"runtime"
)

var currentExecDir string

func init() {
	/* #dogsled */
	_, currentExecFilename, _, _ := runtime.Caller(0)
	currentExecDir = filepath.Dir(currentExecFilename)
}

// GetProjectRootDir gets the project's root directory
// TODO: This assumes a two level deep directory which is very error prone
func GetProjectRootDir() string {
	return currentExecDir + "/../../"
}
