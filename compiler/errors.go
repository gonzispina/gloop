package compiler

import (
	"errors"
	"fmt"
)

func compileErr(t Token, message string) error {
	return errors.New(fmt.Sprintf(
		"Line %v Column %v: %s'",
		t.line,
		t.column,
		message,
	))
}

func unexpectedEndOfFileErr(t Token) error {
	return compileErr(t, "unexpected end of file")
}

func expectedAssignmentOperatorErr(t Token) error {
	return compileErr(t, "expected assignment operator '<-'")
}

func assignErr(t Token, expected varType, got varType) error {
	return compileErr(t, fmt.Sprintf(
		"cannot assign '%s' value to variable of type '%s'",
		got.String(),
		expected.String(),
	))
}

func unexpectedTokenErr(t Token) error {
	return compileErr(t, fmt.Sprintf(
		"unexpected token '%s'",
		t.lexeme,
	))
}

func expectedExpressionErr(t Token) error {
	return compileErr(t, "expected expression")
}

func expectedRightParenthesisErr(t Token) error {
	return compileErr(t, "expected closing parenthesis")
}

func expectedThenErr(t Token) error {
	return compileErr(t, "expected 'then' after expression")
}

func expectedEndIfErr(t Token) error {
	return compileErr(t, "expected 'end if' after block")
}

func expectedTimesErr(t Token) error {
	return compileErr(t, "expected 'times' after expression")
}

func expectedEndLoopErr(t Token) error {
	return compileErr(t, "expected 'end loop' after block")
}

func blockIsTooLargeErr(t Token) error {
	return compileErr(t, "block is too large")
}
