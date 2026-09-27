package core

import "fmt"

var variables = make(map[string]float32)

// SetVariable sets a CDSL variable to the specified value
func SetVariable(variableKey string, variableValue float32) {
	variables[variableKey] = variableValue
}

// HasVariable reports whether a CDSL variable exists
func HasVariable(name string) bool {
	_, ok := variables[name]
	return ok
}

// GetVariable returns the value of a CDSL variable
func GetVariable(name string) float32 {
	return variables[name]
}

// DisplayVariables prints all defined CDSL variables and their values
func DisplayVariables() {
	for key, value := range variables {
		fmt.Printf("%s : %f\n", key, value)
	}
}
