package utils

import (
	"strconv"
	"strings"
)

func Contains(substring string, text string) bool {
	substringLower := strings.ToLower(substring)
	textLower := strings.ToLower(text)

	return strings.Contains(textLower, substringLower)
}

func FormatFloat(value float64, precision int, defaultString string) string {
	if value == 0 {
		return defaultString
	} else {
		return strconv.FormatFloat(value, 'f', precision, 64)
	}
}
