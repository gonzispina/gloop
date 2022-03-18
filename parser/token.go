package parser

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

	LeftArrow
	RightBracket
	LeftBracket
	RightSquareBracket
	LeftSquareBracket
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

	Cell
	Output

	Literal
)

type Token struct {
	tt    tokenType
	value interface{}
}

func literal(value interface{}) Token {
	return Token{tt: Literal, value: value}
}

func token(tt tokenType) Token {
	return Token{tt: tt}
}

func reserved(s string) (Token, error) {
	switch strings.ToLower(s) {
	case "define":
		return token(Define), nil
	case "procedure":
		return token(Procedure), nil
	case "if":
		return token(If), nil
	case "then":
		return token(Then), nil
	case "else":
		return token(Else), nil
	case "and":
		return token(And), nil
	case "or":
		return token(Or), nil
	case "loop":
		return token(Loop), nil
	case "abort":
		return token(Abort), nil
	case "times":
		return token(Times), nil
	case "begin":
		return token(Begin), nil
	case "end":
		return token(End), nil
	case "quit":
		return token(Quit), nil
	case "cell":
		return token(Cell), nil
	case "output":
		return token(Output), nil
	case "yes":
		return literal(true), nil
	case "no":
		return literal(false), nil
	default:
		// Unreachable
		return Token{}, errors.New("not a reserved word")
	}
}
