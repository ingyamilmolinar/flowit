package configs

var supportedFlowitVersions = []string{"0.1"}

// If supported Flowit versions become large we might need to change it for a map
func versionValidator(str string) bool {
	for _, supportedVersion := range supportedFlowitVersions {
		if supportedVersion == str {
			return true
		}
	}
	return false
}
