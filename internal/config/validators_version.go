package config

// If supported Flowit versions become large we might need to change it for a map
func versionValidator(str string) bool {

	var supportedFlowitVersions = []string{"0.1"}

	for _, supportedVersion := range supportedFlowitVersions {
		if supportedVersion == str {
			return true
		}
	}
	return false
}
