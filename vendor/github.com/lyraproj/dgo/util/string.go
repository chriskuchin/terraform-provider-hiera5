package util

import (
	"strconv"
)

// Ftoa returns the given float as a string with almost all trailing zeroes removed. The resulting string will however
// always contain either the letter 'E' or a dot.
func Ftoa(f float64) string {
	s := strconv.FormatFloat(f, 'G', -1, 64)
	for i := range s {
		switch s[i] {
		case 'e', 'E', '.':
			return s
		}
	}
	return s + `.0`
}

// ContainsString returns true if strings contains str
func ContainsString(strings []string, str string) bool {
	for i := range strings {
		if strings[i] == str {
			return true
		}
	}
	return false
}

// StringHash computes a non unique hash code for the given string
func StringHash(s string) int {
	h := 1
	for i := range s {
		h = 31*h + int(s[i])
	}
	return h
}
