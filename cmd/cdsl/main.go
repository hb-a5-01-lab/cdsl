// package main is the entry point
package main

import (
	"fmt"
	"os"

	"github.com/kalexion/cdsl/internal/runtime"
)

func main() {
	if len(os.Args) > 1 {
		if len(os.Args) == 2 && os.Args[1] == "--help" {
			fmt.Println("Usage: cdsl [--help]")
			return
		}

		fmt.Fprintln(os.Stderr, "Error: Unknown option.")
		fmt.Fprintln(os.Stderr, "Usage: cdsl [--help]")
		os.Exit(1)
	}

	// Start the REPL
	runtime.RunInteractiveMode()
}
