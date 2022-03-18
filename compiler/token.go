package compiler

import "errors"

type tokenType uint8

const (
	DoubleQuote tokenType = iota
	Plus
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
	Colon
	Semicolon

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
	Block
	Begin
	End
	Quit

	Cell
	Output
	Yes
	No

	Literal
)

type Token struct {
	tt     tokenType
	lexeme string
	value  interface{}
}

func literal(value interface{}) Token {
	return Token{tt: Literal, value: value}
}

func token(tt tokenType) Token {
	return Token{tt: tt}
}

func reserved(s string) (Token, error) {
	switch s {
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
	case "block":
		return token(Block), nil
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
