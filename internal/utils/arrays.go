package utils

// FindStringInPtrArray returns if array contains the specified string element
func FindStringInPtrArray(s string, arr []*string) bool {
	found := false
	for _, elem := range arr {
		if (*elem) == s {
			found = true
			break
		}
	}
	return found
}

// FindStringInArray returns if array contains the specified string element
func FindStringInArray(s string, arr []string) bool {
	found := false
	for _, elem := range arr {
		if elem == s {
			found = true
			break
		}
	}
	return found
}
