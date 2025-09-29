package utils

import (
	"errors"
	"math"
)

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
		digit := int(s[i] - '0')

		if negative {
			// Check underflow: result*10 - digit >= math.MinInt
			if result < (math.MinInt+digit)/10 {
				return 0, errors.New("integer underflow")
			}
			result = result*10 - digit
		} else {
			// Check overflow: result*10 + digit <= math.MaxInt
			if result > (math.MaxInt-digit)/10 {
				return 0, errors.New("integer overflow")
			}
			result = result*10 + digit
		}
	}

	return result, nil
}
