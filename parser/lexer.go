package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func isNumber(s string) bool {
	return s[0] >= '0' && s[0] <= '9'
}

func isLetter(s string) bool {
	return s[0] >= 'a' && s[0] <= 'z' || s[0] >= 'A' && s[0] <= 'Z' || s[0] == '_'
}

func Lexer(text string) ([]Token, error) {
	var res []Token
	var i int
	text = strings.Trim(text, " ")
	text = strings.Trim(text, "\t")
	text = strings.Trim(text, "\n")

	isAtEnd := func() bool {
		return i >= len(text)
	}

	current := func() string {
		return string(text[i])
	}

	next := func() string {
		l := current()
		i++

		for !isAtEnd() && (l == " " || l == "\n" || l == "\t") {
			l = current()
			i++
		}
		return l
	}

	for !isAtEnd() {
		letter := next()
		switch letter {
		case "\"":
			letter = next()
			for !isAtEnd() && current() != "\"" {
				letter += next()
			}
			if !isAtEnd() {
				next()
			}
			res = append(res, literal(letter))
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
			if !isAtEnd() && current() == "=" {
				res = append(res, token(LesserEqual))
				next()
			} else if current() == "-" {
				res = append(res, token(LeftArrow))
				next()
			} else {
				res = append(res, token(Lesser))
			}
			break
		case ">":
			if current() == "=" {
				res = append(res, token(GreaterEqual))
				next()
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
		case ",":
			res = append(res, token(Comma))
			break
		default:
			if isNumber(letter) {
				lexeme := letter
				for !isAtEnd() && isNumber(current()) {
					lexeme += next()
				}
				value, _ := strconv.ParseInt(lexeme, 10, 64)
				res = append(res, literal(value))
				break
			} else if isLetter(letter) {
				lexeme := letter
				for !isAtEnd() && (isLetter(current()) || isNumber(current())) {
					lexeme += next()
				}
				t, err := reserved(lexeme)
				if err != nil {
					if current() == "?" {
						lexeme += next()
					}
					t = literal(lexeme)
				}
				res = append(res, t)
			} else {
				return nil, errors.New(fmt.Sprintf("unexpected token %s", current()))
			}
		}
	}

	return res, nil
}
