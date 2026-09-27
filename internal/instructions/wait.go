package instructions

import (
	"fmt"
	"strings"
	"time"

	"github.com/znwng/cdsl/internal/core"
)

func waitFunction(delay int) {
	// Placeholder code.
	// Actual delay implementation must eventually be handled by hardware.
	fmt.Printf("Waiting for %d milliseconds\n\n", delay)

	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ProcessWait handles wait commands
func ProcessWait(instruction core.Instruction) {
	if len(instruction) != 2 {
		core.Error(
			"Invalid number of arguments. Example: `WAIT DURATION_MS`",
			instruction,
		)
		return
	}

	value := instruction[1]

	var delay int

	// Expression
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

		delay = int(result)

	case strings.HasPrefix(value, "$"):
		variableName := value[1:]

		if !core.HasVariable(variableName) {
			core.Error("Unknown variable: "+variableName, instruction)
			return
		}

		delay = int(core.GetVariable(variableName))

	default:
		result, err := core.IsValidIntValue(value)
		if err != nil {
			core.Error(err.Error(), instruction)
			return
		}

		delay = result
	}

	// Validate delay
	if delay < 0 {
		core.Error(
			fmt.Sprintf("Delay cannot be negative: %d", delay),
			instruction,
		)
		return
	}

	waitFunction(delay)
}
