package core

import (
	"fmt"
	"strconv"
	"unicode"
)

// IsValidFloatValue parses a string as a 32-bit floating-point value
func IsValidFloatValue(stringValue string) (float32, error) {
	value, err := strconv.ParseFloat(stringValue, 32)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
			return 0, fmt.Errorf("value out of range")
		}

		return 0, fmt.Errorf("invalid value: %s", stringValue)
	}

	return float32(value), nil
}

// IsValidIntValue parses a string as an integer
func IsValidIntValue(stringValue string) (int, error) {
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
			return 0, fmt.Errorf("value out of range")
		}

		return 0, fmt.Errorf("invalid value: %s", stringValue)
	}

	return value, nil
}

// IsValidVariableName checks whether a string is a valid CDSL variable name
func IsValidVariableName(name string) bool {
	if name == "" {
		return false
	}

	runes := []rune(name)

	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}

	for _, char := range runes[1:] {
		if !unicode.IsLetter(char) &&
			!unicode.IsDigit(char) &&
			char != '_' {
			return false
		}
	}

	return true
}
