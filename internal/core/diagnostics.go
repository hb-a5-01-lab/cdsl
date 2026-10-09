package core

import (
	"fmt"
	"os"

	"github.com/kalexion/cdsl/internal/customTypes"
	"github.com/kalexion/cdsl/internal/logging"
	"github.com/kalexion/cdsl/internal/ui"
)

// ========================================
// Instruction Output
// ========================================

// DisplayInstruction prints an instruction and its tokens to standard output.
func DisplayInstruction(instruction customTypes.Instruction) {
	for _, token := range instruction {
		fmt.Printf("%s ", token)
	}

	fmt.Println()
}

// ========================================
// Error Output
// ========================================

// DiagnosticsError prints a formatted error message to standard error
// and logs the failed instruction.
func DiagnosticsError(message string, instruction customTypes.Instruction) {
	fmt.Fprintf(
		os.Stderr,
		"%s%s: %s%s\n\n",
		ui.Red,
		instruction,
		message,
		ui.Reset,
	)

	if err := logging.LogIt(false, message, instruction); err != nil {
		fmt.Fprintf(os.Stderr, "Logging error: %v\n", err)
	}
}

// ========================================
// Success Output
// ========================================

// DiagnosticsSuccess prints a formatted success message to standard output
// and logs the successful instruction.
func DiagnosticsSuccess(message string, instruction customTypes.Instruction) {
	fmt.Printf(
		"%s%s: %s%s\n\n",
		ui.Green,
		instruction,
		message,
		ui.Reset,
	)

	if err := logging.LogIt(true, message, instruction); err != nil {
		fmt.Fprintf(os.Stderr, "Logging error: %v\n", err)
	}
}
