package core

import (
	"fmt"
	"strconv"
	"unicode"
)

// ========================================
// Expression Parsing
// ========================================

// skipWhitespace advances pos past all whitespace characters in expression.
func skipWhitespace(expression string, pos *int) {
	for *pos < len(expression) && unicode.IsSpace(rune(expression[*pos])) {
		*pos++
	}
}

// parseParenthesized parses an expression enclosed in parentheses.
func parseParenthesized(expression string, pos *int) (float32, error) {
	*pos++ // Skip '('

	value, err := parseExpression(expression, pos)
	if err != nil {
		return 0, err
	}

	skipWhitespace(expression, pos)

	if *pos >= len(expression) || expression[*pos] != ')' {
		return 0, fmt.Errorf("expected ')'")
	}

	*pos++ // Skip ')'

	return value, nil
}

// parseVariable parses a variable reference beginning with '$'.
func parseVariable(expression string, pos *int) (float32, error) {
	*pos++ // Skip '$'

	start := *pos

	if *pos >= len(expression) ||
		(!unicode.IsLetter(rune(expression[*pos])) && expression[*pos] != '_') {
		return 0, fmt.Errorf("expected variable name after '$'")
	}

	*pos++

	for *pos < len(expression) &&
		(unicode.IsLetter(rune(expression[*pos])) ||
			unicode.IsDigit(rune(expression[*pos])) ||
			expression[*pos] == '_') {
		*pos++
	}

	name := expression[start:*pos]

	if !HasVariable(name) {
		return 0, fmt.Errorf("undefined variable: $%s", name)
	}

	return GetVariable(name), nil
}

// parseNumber parses a floating-point numeric literal.
func parseNumber(expression string, pos *int) (float32, error) {
	start := *pos

	for *pos < len(expression) &&
		(expression[*pos] >= '0' && expression[*pos] <= '9' ||
			expression[*pos] == '.') {
		(*pos)++
	}

	number := expression[start:*pos]

	value, err := strconv.ParseFloat(number, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", number)
	}

	return float32(value), nil
}

// parsePrimary parses a primary expression such as a number, variable, or
// parenthesized expression.
func parsePrimary(expression string, pos *int) (float32, error) {
	skipWhitespace(expression, pos)

	if *pos >= len(expression) {
		return 0, fmt.Errorf("unexpected end of expression")
	}

	switch expression[*pos] {
	case '(':
		return parseParenthesized(expression, pos)

	case '$':
		return parseVariable(expression, pos)
	}

	if unicode.IsDigit(rune(expression[*pos])) || expression[*pos] == '.' {
		return parseNumber(expression, pos)
	}

	return 0, fmt.Errorf("unexpected character: %c", expression[*pos])
}

// parseFactor parses a factor and handles unary plus and minus operators.
func parseFactor(expression string, pos *int) (float32, error) {
	skipWhitespace(expression, pos)

	if *pos < len(expression) && expression[*pos] == '-' {
		*pos++

		value, err := parseFactor(expression, pos)
		if err != nil {
			return 0, err
		}

		return -value, nil
	}

	if *pos < len(expression) && expression[*pos] == '+' {
		*pos++

		return parseFactor(expression, pos)
	}

	return parsePrimary(expression, pos)
}

// parseTerm parses multiplication and division operations.
func parseTerm(expression string, pos *int) (float32, error) {
	value, err := parseFactor(expression, pos)
	if err != nil {
		return 0, err
	}

	for {
		skipWhitespace(expression, pos)

		if *pos >= len(expression) {
			break
		}

		operation := expression[*pos]

		if operation != '*' && operation != '/' {
			break
		}

		*pos++

		rhs, err := parseFactor(expression, pos)
		if err != nil {
			return 0, err
		}

		if operation == '*' {
			value *= rhs
		} else {
			if rhs == 0 {
				return 0, fmt.Errorf("division by zero")
			}

			value /= rhs
		}
	}

	return value, nil
}

// parseExpression parses addition and subtraction operations.
func parseExpression(expression string, pos *int) (float32, error) {
	value, err := parseTerm(expression, pos)
	if err != nil {
		return 0, err
	}

	for {
		skipWhitespace(expression, pos)

		if *pos >= len(expression) {
			break
		}

		operation := expression[*pos]

		if operation != '+' && operation != '-' {
			break
		}

		*pos++

		rhs, err := parseTerm(expression, pos)
		if err != nil {
			return 0, err
		}

		if operation == '+' {
			value += rhs
		} else {
			value -= rhs
		}
	}

	return value, nil
}

// ========================================
// Expression Evaluation
// ========================================

// EvaluateExpression evaluates an arithmetic expression and returns its
// resulting value.
//
// The expression supports numeric literals, variables, parentheses, unary
// operators, and the standard arithmetic operators +, -, *, and /.
func EvaluateExpression(expression string) (float32, error) {
	pos := 0

	result, err := parseExpression(expression, &pos)
	if err != nil {
		return 0, err
	}

	skipWhitespace(expression, &pos)

	if pos != len(expression) {
		return 0, fmt.Errorf("unexpected character: %c", expression[pos])
	}

	return result, nil
}
