package utils

import (
	"fmt"
	"strconv"
)

func StringToInt(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("error when converting '%s' to int: %w", s, err)
	}
	return num, nil
}

func FloatToInt64(f float64) int64 {
	return int64(f)
}

func StringToFloat64(s string) (float64, error) {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("error when converting '%s' to float64: %v", s, err)
	}
	return result, nil
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
