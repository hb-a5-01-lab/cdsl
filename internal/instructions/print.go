package instructions

import (
	"fmt"
	"strings"

	"github.com/kalexion/cdsl/internal/core"
	"github.com/kalexion/cdsl/internal/customTypes"
)

// ========================================
// Print Processing
// ========================================

// ProcessPrint processes a CDSL PRINT instruction.
//
// The argument may be a variable reference or an arithmetic expression
// enclosed in #[...] syntax.
func ProcessPrint(instruction customTypes.Instruction) {
	if len(instruction) != 2 {
		core.DiagnosticsError(
			"Invalid number of arguments. Example: `PRINT VALUE`",
			instruction,
		)
		return
	}
	argument := instruction[1]

	if strings.HasPrefix(argument, "#[") {
		if len(argument) < 3 || !strings.HasSuffix(argument, "]") {
			core.DiagnosticsError("Invalid expression", instruction)
			return
		}

		expression := argument[2 : len(argument)-1]

		result, err := core.EvaluateExpression(expression)
		if err != nil {
			core.DiagnosticsError(err.Error(), instruction)
			return
		}

		message := fmt.Sprintf("%v", result)
		core.DiagnosticsSuccess(message, instruction)
		return
	}

	variableKey := strings.TrimPrefix(argument, "$")

	if !core.HasVariable(variableKey) {
		core.DiagnosticsError(
			"No variable with name "+variableKey,
			instruction,
		)
		return
	}

	message := fmt.Sprintf("%v", core.GetVariable(variableKey))
	core.DiagnosticsSuccess(message, instruction)
}
