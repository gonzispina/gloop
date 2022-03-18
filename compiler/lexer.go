package compiler

import (
	"errors"
	"fmt"
	"strconv"
)

func isNumber(s uint8) bool {
	return s >= '0' && s <= '9'
}

func isLetter(s uint8) bool {
	return s >= 'a' && s <= 'z' || s >= 'A' && s <= 'Z' || s == '_'
}

func Lexer(text string) ([]Token, error) {
	var res []Token
	var i int

	isAtEnd := func() bool {
		return i == len(text)-1
	}

	letter := func() string {
		l := string(text[i])
		for !isAtEnd() && (l == " " || l == "\n") {
			l = string(text[i])
		}
		return l
	}

	nextLetter := func() string {
		l := string(text[i+1])
		for !isAtEnd() && (l == " " || l == "\n") {
			l = string(text[i+1])
		}
		return l
	}

	for i = 0; i < len(text); i++ {
		switch letter() {
		case " ":
			continue
		case "\n":
			continue
		case "\"":
			lexeme := ""
			for !isAtEnd() && letter() != "\"" {
				lexeme += letter()
				i++
			}
			res = append(res, literal(lexeme))
			break
		case "+":
			res = append(res, token(Plus))
			break
		case "*":
			res = append(res, token(Star))
			break
		case "=":
			res = append(res, token(Equal))
			break
		case "<":
			if !isAtEnd() && nextLetter() == "=" {
				res = append(res, token(LesserEqual))
			} else if nextLetter() == "-" {
				res = append(res, token(LeftArrow))
			} else {
				res = append(res, token(Lesser))
			}
			break
		case ">":
			if nextLetter() == "=" {
				res = append(res, token(GreaterEqual))
			} else {
				res = append(res, token(Greater))
			}
			break
		case "(":
			res = append(res, token(LeftBracket))
			break
		case ")":
			res = append(res, token(RightBracket))
			break
		case "[":
			res = append(res, token(LeftSquareBracket))
			break
		case "]":
			res = append(res, token(RightSquareBracket))
			break
		case ":":
			res = append(res, token(Colon))
			break
		case ";":
			res = append(res, token(Semicolon))
			break
		default:
			lexeme := letter()
			if isNumber(text[i]) {
				for !isAtEnd() && isNumber(text[i]) {
					lexeme += letter()
					i++
				}
				value, _ := strconv.ParseInt(lexeme, 10, 64)
				res = append(res, literal(value))
				break
			} else if isLetter(text[i]) {
				for !isAtEnd() && (isLetter(text[i]) || isNumber(text[i])) {
					lexeme += letter()
					i++
				}
				t, err := reserved(lexeme)
				if err != nil {
					t = literal(lexeme)
				}
				res = append(res, t)
			} else {
				return nil, errors.New(fmt.Sprintf("unexpected token %s", letter()))
			}
		}
	}

	return res, nil
}
