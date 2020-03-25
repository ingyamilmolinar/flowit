package utils

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// GetRootDirectory gets the project root directory
func GetRootDirectory() string {
	return basepath + "/../../"
}
