package util

import "strings"

// Contains to check if the slice has an element specified by the keyword
func Contains(slice []string, keyword string) bool {
	if len(slice) == 0 {
		return false
	}

	for _, ele := range slice {
		if ele == keyword {
			return true
		}
	}
	return false
}

// ContainsIgnoreCase to check if the slice has an element specified by the keyword, case insensitive
func ContainsIgnoreCase(slice []string, keyword string) bool {
	if len(slice) == 0 {
		return false
	}

	for _, ele := range slice {
		if strings.ToLower(ele) == strings.ToLower(keyword) {
			return true
		}
	}
	return false
}

// Unique to remove the duplicated element from a string slice
func Unique(slice []string) []string {
	var results []string
	m := make(map[string]struct{})

	for _, ele := range slice {
		if _, ok := m[ele]; !ok {
			results = append(results, ele)
			m[ele] = struct{}{}
		}
	}

	return results
}

// SliceTrim to remove the element which equal space
func SliceTrim(slice []string) []string {
	results := []string{}
	for _, ele := range slice {
		str := strings.TrimSpace(ele)
		if str != "" {
			results = append(results, str)
		}
	}

	return results
}

// Remove to remove the element with value 'v' from slice
func Remove(slice []string, v string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			slice = append(slice[:i], slice[i+1:]...)
			i-- // maintain the correct index
		}
	}
	return slice
}
