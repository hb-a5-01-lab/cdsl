package core

import (
	"fmt"
	"strconv"
	"unicode"
)

// ========================================
// Value Validation
// ========================================

// IsValidFloatValue parses a string as a 32-bit floating-point value.
//
// An error is returned when the string is not a valid floating-point value or
// when the value cannot be represented as a float32.
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

// IsValidIntValue parses a string as an integer.
//
// An error is returned when the string is not a valid integer or when the
// value cannot be represented as an int.
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

// ========================================
// Variable Validation
// ========================================

// IsValidVariableName reports whether name is a valid CDSL variable name.
//
// A variable name must begin with a letter or underscore and may contain
// letters, digits, and underscores thereafter.
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
