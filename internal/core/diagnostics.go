package core

import (
	"fmt"
	"os"

	"github.com/znwng/cdsl/internal/ui"
)

// DisplayInstruction prints the instruction
func DisplayInstruction(instruction Instruction) {
	for _, token := range instruction {
		fmt.Printf("%s ", token)
	}

	fmt.Println()
}

// Error is used to print a formatted colored error
func Error(message string, instruction Instruction) {
	fmt.Fprintf(os.Stderr, "%s%s: %s%s\n", ui.Red, instruction, message, ui.Reset)
}
