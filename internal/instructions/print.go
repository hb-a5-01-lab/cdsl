package instructions

import (
	"fmt"
	"strings"

	"github.com/znwng/cdsl/internal/core"
)

// ========================================
// Print Processing
// ========================================

// ProcessPrint processes a CDSL PRINT instruction.
//
// The argument may be a variable reference or an arithmetic expression
// enclosed in #[...] syntax.
func ProcessPrint(instruction core.Instruction) {
	if len(instruction) != 2 {
		core.Error(
			"Invalid number of arguments. Example: `PRINT VALUE`",
			instruction,
		)
		return
	}

	argument := instruction[1]

	if strings.HasPrefix(argument, "#[") {
		if len(argument) < 3 || !strings.HasSuffix(argument, "]") {
			core.Error("Invalid expression", instruction)
			return
		}

		expression := argument[2 : len(argument)-1]

		result, err := core.EvaluateExpression(expression)
		if err != nil {
			core.Error(err.Error(), instruction)
			return
		}

		fmt.Println(result)
		return
	}

	variableKey := strings.TrimPrefix(argument, "$")

	if !core.HasVariable(variableKey) {
		core.Error(
			"No variable with name "+variableKey,
			instruction,
		)
		return
	}

	fmt.Println(core.GetVariable(variableKey))
}
