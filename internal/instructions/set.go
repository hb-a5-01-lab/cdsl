package instructions

import (
	"fmt"
	"strings"

	"github.com/znwng/cdsl/internal/core"
)

// ProcessSet sets the variables and handles them
func ProcessSet(instruction core.Instruction) {
	if len(instruction) != 3 {
		core.Error(
			"Invalid number of arguments. Example: `SET VARIABLE VALUE`",
			instruction,
		)
		return
	}

	if !core.IsValidVariableName(instruction[1]) {
		core.Error(
			"Invalid variable name: "+instruction[1],
			instruction,
		)
		return
	}

	variableKey := instruction[1]
	value := instruction[2]

	var variableValue float32

	switch {
	case strings.HasPrefix(value, "#["):
		if len(value) < 3 || !strings.HasSuffix(value, "]") {
			core.Error("Invalid expression", instruction)
			return
		}

		expression := value[2 : len(value)-1]

		result, err := core.EvaluateExpression(expression)
		if err != nil {
			core.Error(err.Error(), instruction)
			return
		}

		variableValue = result

	case strings.HasPrefix(value, "$"):
		variableName := value[1:]

		if !core.HasVariable(variableName) {
			core.Error("Unknown variable: "+variableName, instruction)
			return
		}

		variableValue = core.GetVariable(variableName)

	default:
		result, err := core.IsValidFloatValue(value)
		if err != nil {
			core.Error(err.Error(), instruction)
			return
		}

		variableValue = result
	}

	// Store variable
	core.SetVariable(variableKey, variableValue)

	fmt.Printf("variable %s set to %v\n", variableKey, variableValue)
}
