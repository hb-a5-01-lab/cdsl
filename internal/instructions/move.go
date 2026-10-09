// Package instructions defines all instructions
package instructions

import (
	"fmt"
	"strings"

	"github.com/kalexion/cdsl/internal/core"
	"github.com/kalexion/cdsl/internal/customTypes"
	"github.com/kalexion/cdsl/internal/hardware"
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
	instruction customTypes.Instruction,
	componentLabel string,
	valueText string,
	value float32,
) bool {
	withinLimits, err := core.ValueWithinLimits(config, componentLabel, value)

	if err != nil {
		core.DiagnosticsError(err.Error(), instruction)
		return false
	}

	if !withinLimits {
		core.DiagnosticsError(
			fmt.Sprintf(
				"Value out of range for %s: %f",
				componentLabel,
				value,
			),
			instruction,
		)
		return false
	}

	if err := hardware.SendCommand(config, componentLabel, value); err != nil {
		core.DiagnosticsError(err.Error(), instruction)
		return false
	}

	core.DiagnosticsSuccess(
		fmt.Sprintf("Moved %s by %s", componentLabel, valueText),
		instruction,
	)

	return true

}

// ========================================
// Move Processing
// ========================================

// ProcessMove processes a CDSL MOVE instruction.
//
// The movement value may be a numeric value, a variable reference, or an
// arithmetic expression enclosed in #[...] syntax.
func ProcessMove(config *core.Config, instruction customTypes.Instruction) {
	if len(instruction) != 3 {
		core.DiagnosticsError(
			"Invalid number of arguments. "+
				"Example: `MOVE COMPONENT_NAME VALUE`",
			instruction,
		)
		return
	}

	componentLabel := instruction[1]

	componentExists, err := core.ComponentExists(config, componentLabel)

	if err != nil {
		core.DiagnosticsError(err.Error(), instruction)
		return
	}

	if !componentExists {
		core.DiagnosticsError("Undefined component: "+componentLabel, instruction)
		return
	}

	value := instruction[2]

	if strings.HasPrefix(value, "#[") {
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
			core.DiagnosticsError("Unknown variable: "+variableName, instruction)
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
		core.DiagnosticsError(err.Error(), instruction)
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
