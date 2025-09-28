package utils

import "errors"

func FastStringToInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("empty string is not a valid integer")
	}

	var negative bool
	var start int
	if s[0] == '-' {
		if len(s) == 1 {
			return 0, errors.New("invalid integer string")
		}
		negative = true
		start = 1
	}
	
	var result int
	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, errors.New("invalid integer string")
		}
		result = result*10 + int(s[i]-'0')
	}
	if negative {
		result = -result
	}

	return result, nil
}
