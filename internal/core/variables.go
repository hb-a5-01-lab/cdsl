package core

import "fmt"

var variables = make(map[string]float32)

// ========================================
// Variable Storage
// ========================================

// SetVariable assigns a value to a CDSL variable.
//
// If the variable already exists, its current value is replaced.
func SetVariable(variableKey string, variableValue float32) {
	variables[variableKey] = variableValue
}

// HasVariable reports whether a CDSL variable is defined.
func HasVariable(name string) bool {
	_, ok := variables[name]
	return ok
}

// GetVariable returns the value of a CDSL variable.
//
// If the variable is not defined, its zero value is returned.
func GetVariable(name string) float32 {
	return variables[name]
}

// ========================================
// Variable Output
// ========================================

// DisplayVariables prints all defined CDSL variables and their values.
func DisplayVariables() {
	for key, value := range variables {
		fmt.Printf("%s : %f\n", key, value)
	}
}
