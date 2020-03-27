package utils

// ExitIfErr logs and panics if an error exists
func ExitIfErr(err error) {
	if err != nil {
		logger := GetLogger()
		logger.Error(err)
		panic(err)
	}
}
