package instructions

import (
	"fmt"
	"strings"
	"time"

	"github.com/kalexion/cdsl/internal/core"
	"github.com/kalexion/cdsl/internal/customTypes"
)

// ========================================
// Wait Execution
// ========================================

// waitFunction pauses execution for the specified duration in milliseconds.
//
// The delay is currently handled by the host system and may be moved to the
// hardware layer in the future.
func waitFunction(delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ========================================
// Wait Processing
// ========================================

// ProcessWait processes a CDSL WAIT instruction.
//
// The duration may be an integer value, a variable reference, or an arithmetic
// expression enclosed in #[...] syntax. Negative durations are rejected.
func ProcessWait(instruction customTypes.Instruction) {
	if len(instruction) != 2 {
		core.DiagnosticsError(
			"Invalid number of arguments. Example: `WAIT DURATION_MS`",
			instruction,
		)
		return
	}

	value := instruction[1]

	var delay int

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

		delay = int(result)

	case strings.HasPrefix(value, "$"):
		variableName := value[1:]

		if !core.HasVariable(variableName) {
			core.DiagnosticsError("Unknown variable: "+variableName, instruction)
			return
		}

		delay = int(core.GetVariable(variableName))

	default:
		result, err := core.IsValidIntValue(value)
		if err != nil {
			core.DiagnosticsError(err.Error(), instruction)
			return
		}

		delay = result
	}

	if delay < 0 {
		core.DiagnosticsError(
			fmt.Sprintf("Delay cannot be negative: %d", delay),
			instruction,
		)
		return
	}

	core.DiagnosticsSuccess(
		fmt.Sprintf("Waiting for %d milliseconds", delay),
		instruction,
	)

	waitFunction(delay)
}
