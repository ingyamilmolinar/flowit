package io

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetProjectRootDir gets the project's root directory
// TODO: This assumes a two level deep directory which is very error prone
func GetProjectRootDir() string {
	/* #dogsled */
	_, currentExecFilename, _, _ := runtime.Caller(0)
	return filepath.Dir(currentExecFilename) + "/../../"
}

// RemoveDirectory receives a string that represents a filesystem directory path and deletes it and all it's children.
// An error is returned in case of failure
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}
