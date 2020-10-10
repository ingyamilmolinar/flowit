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

func MergeSlices(a, b []string) []string {
	var result []string
	sizeA := len(a)
	sizeB := len(b)
	maxSize := maxInt(sizeA, sizeB)
	for i := 0; i < maxSize; i++ {
		if i < sizeA {
			result = append(result, a[i])
		}
		if i < sizeB {
			result = append(result, b[i])
		}
	}
	return result
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
