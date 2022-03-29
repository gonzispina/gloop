package compiler

import (
	"errors"
	"strings"
)

type tokenType uint8

const (
	Plus tokenType = iota
	Star
	Equal
	Lesser
	LesserEqual
	Greater
	GreaterEqual
	Not

	LeftArrow
	LeftParen
	RightParen
	LeftSquareBracket
	RightSquareBracket
	Comma

	Define
	Procedure
	If
	Then
	Else
	And
	Or
	Loop
	Abort
	Times
	Begin
	End
	Quit

	Output

	Identifier
	Constant

	Eof
)

type Token struct {
	tt     tokenType
	lexeme string
	value  interface{}
	line   int
	column int
}

func identifier(lexeme string, value interface{}, line, column int) Token {
	return Token{tt: Identifier, lexeme: lexeme, value: value, line: line, column: column}
}

func constant(lexeme string, value interface{}, line, column int) Token {
	return Token{tt: Constant, lexeme: lexeme, value: value, line: line, column: column}
}

func token(tt tokenType, lexeme string, line, column int) Token {
	return Token{tt: tt, lexeme: lexeme, line: line, column: column}
}

func reserved(s string, line, column int) (Token, error) {
	switch strings.ToLower(s) {
	case "define":
		return token(Define, strings.ToUpper(s), line, column), nil
	case "procedure":
		return token(Procedure, strings.ToUpper(s), line, column), nil
	case "if":
		return token(If, strings.ToUpper(s), line, column), nil
	case "then":
		return token(Then, strings.ToUpper(s), line, column), nil
	case "else":
		return token(Else, strings.ToUpper(s), line, column), nil
	case "not":
		return token(Not, strings.ToUpper(s), line, column), nil
	case "and":
		return token(And, strings.ToUpper(s), line, column), nil
	case "or":
		return token(Or, strings.ToUpper(s), line, column), nil
	case "loop":
		return token(Loop, strings.ToUpper(s), line, column), nil
	case "abort":
		return token(Abort, strings.ToUpper(s), line, column), nil
	case "times":
		return token(Times, strings.ToUpper(s), line, column), nil
	case "begin":
		return token(Begin, strings.ToUpper(s), line, column), nil
	case "end":
		return token(End, strings.ToUpper(s), line, column), nil
	case "quit":
		return token(Quit, strings.ToUpper(s), line, column), nil
	case "output":
		return token(Output, strings.ToUpper(s), line, column), nil
	case "yes":
		return constant(strings.ToUpper(s), true, line, column), nil
	case "no":
		return constant(strings.ToUpper(s), false, line, column), nil
	default:
		// Unreachable
		return Token{}, errors.New("not a reserved word")
	}
}
