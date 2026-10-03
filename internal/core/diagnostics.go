package core

import (
	"fmt"
	"os"

	"github.com/znwng/cdsl/internal/ui"
)

// ========================================
// Instruction Output
// ========================================

// DisplayInstruction prints an instruction and its tokens to standard output.
func DisplayInstruction(instruction Instruction) {
	for _, token := range instruction {
		fmt.Printf("%s ", token)
	}

	fmt.Println()
}

// ========================================
// Error Output
// ========================================

// Error prints a formatted error message to standard error.
//
// The instruction that caused the error is included in the message to provide
// context for the failure.
func Error(message string, instruction Instruction) {
	fmt.Fprintf(os.Stderr, "%s%s: %s%s\n\n", ui.Red, instruction, message, ui.Reset)
}
