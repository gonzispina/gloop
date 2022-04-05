package compiler

import (
	"errors"
	"fmt"
)

type ErrCode string

const (
	UnexpectedEOFErrCode              = "Unexpected end of file"
	ExpectedAssignmentOperatorErrCode = "Expected assignment operator"
	InvalidTypeErrCode                = "Invalid variable type"
	UnexpectedTokenErrCode            = "Unexpected token"
	ExpectedExpressionErrCode         = "Expected expression"
	ExpectedRightParenErrCode         = "Expected right parenthesis"
	ExpectedThenErrCode               = "Expected then"
	ExpectedEndIfErrCode              = "Expected end if"
	ExpectedTimesAfterLoopErrCode     = "Expected times"
	ExpectedEndLoopAfterLoopErrCode   = "Expected end loop"
	UndefinedVariableErrCode          = "Undefined variable"
	BlockIsTooLargeErrCode            = "Block is too large"
	BooleanExpressionNeededCodeErr    = "Boolean expression needed"
	NumberExpressionNeededCodeErr     = "Number expression needed"
)

func compileErr(t Token, message string, code ErrCode) error {
	return errors.New(fmt.Sprintf(
		"Line %v Column %v: %s'",
		t.line,
		t.column,
		fmt.Sprintf("%s. ErrCode: %s", message, string(code)),
	))
}

func unexpectedEndOfFileErr(t Token) error {
	return compileErr(t, "unexpected end of file", UnexpectedEOFErrCode)
}

func expectedAssignmentOperatorErr(t Token) error {
	return compileErr(t, "expected assignment operator '<-'", ExpectedAssignmentOperatorErrCode)
}

func invalidTypeErr(t Token, expected varType, got varType) error {
	return compileErr(t, fmt.Sprintf(
		"cannot assign '%s' value to variable of type '%s'",
		got.String(),
		expected.String(),
	), InvalidTypeErrCode)
}

func unexpectedTokenErr(t Token) error {
	return compileErr(t, fmt.Sprintf(
		"unexpected token '%s'",
		t.lexeme,
	), UnexpectedTokenErrCode)
}

func expectedExpressionErr(t Token) error {
	return compileErr(t, "expected expression", ExpectedExpressionErrCode)
}

func expectedRightParenthesisErr(t Token) error {
	return compileErr(t, "expected closing parenthesis", ExpectedRightParenErrCode)
}

func expectedThenErr(t Token) error {
	return compileErr(t, "expected 'then' after expression", ExpectedThenErrCode)
}

func expectedEndIfErr(t Token) error {
	return compileErr(t, "expected 'end if' after block", ExpectedEndIfErrCode)
}

func expectedTimesErr(t Token) error {
	return compileErr(t, "expected 'times' after expression", ExpectedTimesAfterLoopErrCode)
}

func expectedEndLoopErr(t Token) error {
	return compileErr(t, "expected 'end loop' after block", ExpectedEndLoopAfterLoopErrCode)
}

func blockIsTooLargeErr(t Token) error {
	return compileErr(t, "block is too large", BlockIsTooLargeErrCode)
}

func undefinedVariableErr(t Token, name string) error {
	return compileErr(t,
		fmt.Sprintf("cannot evaluate '%s' because it wasn't assigned before: %s <- Value",
			name,
			name,
		),
		UndefinedVariableErrCode,
	)
}

func booleanExpressionNeededErr(t Token) error {
	return compileErr(
		t,
		"needed boolean expression",
		BooleanExpressionNeededCodeErr,
	)
}

func numberExpressionNeededErr(t Token) error {
	return compileErr(
		t,
		"needed number expression",
		NumberExpressionNeededCodeErr,
	)
}
