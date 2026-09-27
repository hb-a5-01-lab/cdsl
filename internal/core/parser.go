package core

import "strings"

// Tokenize takes the instruction as string and returns the Instruction type
func Tokenize(line string) Instruction {
	var tokens Instruction

	for token := range strings.FieldsSeq(line) {
		if strings.HasPrefix(token, "//") {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}
