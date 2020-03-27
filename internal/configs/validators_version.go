package configs

var supportedFlowitVersions = []string{"0.1"}

func versionValidator(str string) bool {
	for _, supportedVersion := range supportedFlowitVersions {
		if supportedVersion == str {
			return true
		}
	}
	return false
}
