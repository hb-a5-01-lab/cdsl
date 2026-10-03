// Package instructions defines all instructions
package instructions

import (
	"fmt"
	"strings"

	"github.com/znwng/cdsl/internal/core"
	"github.com/znwng/cdsl/internal/hardware"
)

// ========================================
// Move Execution
// ========================================

// ExecuteMove validates and executes a component movement command.
//
// The requested value is checked against the component's configured limits
// before the command is sent to the Arduino.
func ExecuteMove(
	config *core.Config,
	instruction core.Instruction,
	componentLabel string,
	valueText string,
	value float32,
) bool {
	withinLimits, err := core.ValueWithinLimits(config, componentLabel, value)

	if err != nil {
		core.Error(err.Error(), instruction)
		return false
	}

	if !withinLimits {
		core.Error(
			fmt.Sprintf(
				"Value out of range for %s: %f",
				componentLabel,
				value,
			),
			instruction,
		)
		return false
	}

	err = hardware.SendCommand(config, componentLabel, value)

	if err != nil {
		core.Error(err.Error(), instruction)
		return false
	}

	fmt.Printf("Moved %s by %s\n\n", componentLabel, valueText)

	return true
}

// ========================================
// Move Processing
// ========================================

// ProcessMove processes a CDSL MOVE instruction.
//
// The movement value may be a numeric value, a variable reference, or an
// arithmetic expression enclosed in #[...] syntax.
func ProcessMove(config *core.Config, instruction core.Instruction) {
	if len(instruction) != 3 {
		core.Error(
			"Invalid number of arguments. "+
				"Example: `MOVE COMPONENT_NAME VALUE`",
			instruction,
		)
		return
	}

	componentLabel := instruction[1]

	componentExists, err := core.ComponentExists(config, componentLabel)

	if err != nil {
		core.Error(err.Error(), instruction)
		return
	}

	if !componentExists {
		core.Error("Undefined component: "+componentLabel, instruction)
		return
	}

	value := instruction[2]

	if strings.HasPrefix(value, "#[") {
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

		ExecuteMove(
			config,
			instruction,
			componentLabel,
			value,
			result,
		)

		return
	}

	if strings.HasPrefix(value, "$") {
		variableName := value[1:]

		if !core.HasVariable(variableName) {
			core.Error(
				"Unknown variable: "+variableName,
				instruction,
			)
			return
		}

		result := core.GetVariable(variableName)

		ExecuteMove(
			config,
			instruction,
			componentLabel,
			value,
			result,
		)

		return
	}

	result, err := core.IsValidFloatValue(value)

	if err != nil {
		core.Error(err.Error(), instruction)
		return
	}

	ExecuteMove(
		config,
		instruction,
		componentLabel,
		value,
		result,
	)
}
