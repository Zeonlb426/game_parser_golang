package utils

import (
	"strings"
)

func IntInSlice(int int, list []int) bool {
	for _, a := range list {
		if a == int {
			return true
		}
	}

	return false
}

func Int8InSlice(int int8, list []int8) bool {
	for _, a := range list {
		if a == int {
			return true
		}
	}

	return false
}

func StringInSlice(string string, list []string) bool {
	for _, a := range list {
		if a == string {
			return true
		}
	}

	return false
}

func ContainsInSlice(substrings []string, text string) bool {
	textLower := strings.ToLower(text)

	for _, substring := range substrings {
		substringLower := strings.ToLower(substring)

		if strings.Contains(textLower, substringLower) {
			return true
		}
	}

	return false
}
