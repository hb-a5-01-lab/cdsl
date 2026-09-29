package runtime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/znwng/cdsl/internal/core"
	"github.com/znwng/cdsl/internal/instructions"
	"github.com/znwng/cdsl/internal/ui"

	"github.com/chzyer/readline"
)

// ========================================
// Instruction Processing
// ========================================

type opcode int

const (
	set opcode = iota
	print
	move
	wait
	printc
	reset
	invalid
)

var opcodeMap = map[string]opcode{
	"set":    set,
	"print":  print,
	"move":   move,
	"wait":   wait,
	"printc": printc,
	"reset":  reset,
}

func processInstruction(config *core.Config, instruction core.Instruction) {
	if len(instruction) == 0 {
		return
	}

	action, ok := opcodeMap[strings.ToLower(instruction[0])]
	if !ok {
		core.Error(
			"Invalid action: "+instruction[0],
			instruction,
		)
		return
	}

	switch action {
	case set:
		instructions.ProcessSet(instruction)

	case print:
		instructions.ProcessPrint(instruction)

	case move:
		instructions.ProcessMove(config, instruction)

	case wait:
		instructions.ProcessWait(instruction)

	case printc:
		instructions.ProcessPrintc(config)

	case reset:
		instructions.ProcessReset(config, instruction)

	default:
		core.Error("Invalid action: "+instruction[0], instruction)
	}

	fmt.Println()
}

// ========================================
// Interactive Mode
// ========================================

// RunInteractiveMode starts the CDSL interactive REPL.
//
// The REPL loads the user's configuration and maintains command history in
// ~/.cdsl_history.
func RunInteractiveMode() {
	config, err := core.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	historyFile := filepath.Join(home, ".cdsl_history")

	fmt.Printf("%sCDSL Interactive Mode%s\n", ui.Green, ui.Reset)
	fmt.Println("Type `exit` to quit.")
	fmt.Println("Type `clear` to clear screen.")
	fmt.Println()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ui.Green + "cdsl> " + ui.Reset,
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	defer func() {
		if err := rl.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			}

			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line == "exit" {
			break
		}

		if line == "clear" {
			fmt.Print("\033[2J\033[H")
			continue
		}

		instruction := core.Tokenize(line)

		if len(instruction) == 0 {
			continue
		}

		processInstruction(config, instruction)
	}
}
