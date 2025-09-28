package utils

import "errors"

func FastStringToInt(s string) (int, error) {
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, errors.New("invalid integer string")
		}
		result = result*10 + int(s[i]-'0')
	}

	return result, nil
}
