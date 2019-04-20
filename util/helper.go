package util

import (
	"strings"
)

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Uint8 is a helper routine that allocates a new uint8 value
// to store v and returns a pointer to it.
func Uint8(v uint8) *uint8 { return &v }

// Uint64 is a helper routine that allocates a new uint64 value
// to store v and returns a pointer to it.
func Uint64(v uint64) *uint64 { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// SafeString return the value of string pointer
func SafeString(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func SafeStringSlice(v *[]string) []string {
	if v != nil {
		return *v
	}
	return []string{}
}

// SafeToLower change string to lower case
func SafeToLower(v *string) *string {
	if v != nil {
		return String(strings.ToLower(*v))
	}
	return v
}

// SafeToUpper change string to upper case
func SafeToUpper(v *string) *string {
	if v != nil {
		return String(strings.ToUpper(*v))
	}
	return v
}

// IsTrue returns true or false according to the pointer to bool type.
func IsTrue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// EqualsIgnoreCase to check if the two string equals, case insensitive
func EqualsIgnoreCase(str1, str2 string) bool {
	return strings.ToLower(str1) == strings.ToLower(str2)
}
