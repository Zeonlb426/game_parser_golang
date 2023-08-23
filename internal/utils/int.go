package utils

func BoolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func IntToBool(value int) bool {
	if value > 1 {
		return true
	}

	return false
}
