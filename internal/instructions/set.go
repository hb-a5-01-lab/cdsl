package instructions

import (
	"fmt"
	"strings"

	"github.com/kalexion/cdsl/internal/core"
	"github.com/kalexion/cdsl/internal/customTypes"
)

// ========================================
// Variable Processing
// ========================================

// ProcessSet processes a CDSL SET instruction.
//
// The assigned value may be a numeric value, a variable reference, or an
// arithmetic expression enclosed in #[...] syntax.
func ProcessSet(instruction customTypes.Instruction) {
	if len(instruction) != 3 {
		core.DiagnosticsError(
			"Invalid number of arguments. Example: `SET VARIABLE VALUE`",
			instruction,
		)
		return
	}

	if !core.IsValidVariableName(instruction[1]) {
		core.DiagnosticsError(
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
			core.DiagnosticsError("Invalid expression", instruction)
			return
		}

		expression := value[2 : len(value)-1]

		result, err := core.EvaluateExpression(expression)
		if err != nil {
			core.DiagnosticsError(err.Error(), instruction)
			return
		}

		variableValue = result

	case strings.HasPrefix(value, "$"):
		variableName := value[1:]

		if !core.HasVariable(variableName) {
			core.DiagnosticsError("Unknown variable: "+variableName, instruction)
			return
		}

		variableValue = core.GetVariable(variableName)

	default:
		result, err := core.IsValidFloatValue(value)
		if err != nil {
			core.DiagnosticsError(err.Error(), instruction)
			return
		}

		variableValue = result
	}

	core.SetVariable(variableKey, variableValue)
	core.DiagnosticsSuccess(fmt.Sprintf("Variable %s set to %v", variableKey, variableValue), instruction)
}
