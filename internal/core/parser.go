package core

import (
	"strings"

	"github.com/kalexion/cdsl/internal/customTypes"
)

// Tokenize splits an instruction line into individual tokens.
//
// Tokenization stops when a comment marker ("//") is encountered. The comment
// and any tokens following it are excluded from the resulting instruction.
func Tokenize(line string) customTypes.Instruction {
	var tokens customTypes.Instruction

	for token := range strings.FieldsSeq(line) {
		if strings.HasPrefix(token, "//") {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}
