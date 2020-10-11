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

// MergeSlices takes two string lists and merge them into a single list
// TODO: Unit test
func MergeSlices(a, b []string) []string {
	var result []string
	sizeA := len(a)
	sizeB := len(b)
	maxSize := MaxInt(sizeA, sizeB)
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

// MaxInt returns the maximum of two integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CompareSlices returns true if the slices are equal length with equal elements
// TODO: Unit test
func CompareSlices(a, b []string) bool {
	sizeA := len(a)
	sizeB := len(b)
	if sizeA != sizeB {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
