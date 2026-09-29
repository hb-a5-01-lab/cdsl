package core

import "strings"

// Tokenize splits an instruction line into individual tokens.
//
// Tokenization stops when a comment marker ("//") is encountered. The comment
// and any tokens following it are excluded from the resulting instruction.
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
