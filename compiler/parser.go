package compiler

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
	var i, line int

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
			if l == "\n" {
				line++
			}

			l = current()
			i++
		}
		return l
	}

	for !isAtEnd() {
		letter := next()
		switch letter {
		case "\"":
			lexeme := next()
			for !isAtEnd() && current() != "\"" {
				lexeme += next()
			}
			if !isAtEnd() {
				next()
			}
			res = append(res, identifier(letter, lexeme, line, i))
			break
		case "+":
			res = append(res, token(Plus, letter, line, i))
			break
		case "*":
			res = append(res, token(Star, letter, line, i))
			break
		case "=":
			res = append(res, token(Equal, letter, line, i))
			break
		case "<":
			if !isAtEnd() && current() == "=" {
				res = append(res, token(LesserEqual, letter+next(), line, i))
			} else if current() == "-" {
				res = append(res, token(LeftArrow, letter+next(), line, i))
			} else {
				res = append(res, token(Lesser, letter, line, i))
			}
			break
		case ">":
			if current() == "=" {
				res = append(res, token(GreaterEqual, letter+next(), line, i))
			} else {
				res = append(res, token(Greater, letter, line, i))
			}
			break
		case "(":
			res = append(res, token(LeftParen, letter, line, i))
			break
		case ")":
			res = append(res, token(RightParen, letter, line, i))
			break
		case "[":
			res = append(res, token(LeftSquareBracket, letter, line, i))
			break
		case "]":
			res = append(res, token(RightSquareBracket, letter, line, i))
			break
		case ",":
			res = append(res, token(Comma, letter, line, i))
			break
		default:
			if isNumber(letter) {
				lexeme := letter
				for !isAtEnd() && isNumber(current()) {
					lexeme += next()
				}
				value, _ := strconv.ParseInt(lexeme, 10, 64)
				res = append(res, constant(lexeme, value, line, i))
				break
			} else if isLetter(letter) {
				lexeme := letter
				for !isAtEnd() && (isLetter(current()) || isNumber(current())) {
					lexeme += next()
				}

				if strings.ToLower(lexeme) == "end" {
					lexeme += next()
					for !isAtEnd() && isLetter(current()) {
						lexeme += next()
					}

					switch strings.ToLower(lexeme) {
					case "endif":
					case "endloop":
					case "endprocedure":
						break
					default:
						return nil, errors.New(fmt.Sprintf("Line %v Column %v: unexpected token %s. Expected 'end procedure', 'end if' or 'end loop' statement", line, i, current()))
					}
				}

				if strings.ToLower(lexeme) == "define" {
					lexeme += next()
					for !isAtEnd() && isLetter(current()) {
						lexeme += next()
					}

					if strings.ToLower(lexeme) != "defineprocedure" {
						return nil, errors.New(fmt.Sprintf("Line %v Column %v: unexpected token %s. Expected 'define procedure' statement", line, i, current()))
					}
				}

				if strings.ToLower(lexeme) == "quit" {
					lexeme += next()
					for !isAtEnd() && isLetter(current()) {
						lexeme += next()
					}

					if strings.ToLower(lexeme) != "quitprocedure" {
						return nil, errors.New(fmt.Sprintf("Line %v Column %v: unexpected token %s. Expected 'quit procedure' statement", line, i, current()))
					}
				}

				if strings.ToLower(lexeme) == "abort" {
					lexeme += next()
					for !isAtEnd() && (isLetter(current()) || isNumber(current())) {
						lexeme += next()
					}

					if strings.ToLower(lexeme) != "abortloop" {
						return nil, errors.New(fmt.Sprintf("Line %v Column %v: unexpected token %s. Expected 'abort loop' statement", line, i, current()))
					}
				}

				t, err := reserved(lexeme, line, i)
				if err != nil {
					if current() == "?" {
						lexeme += next()
					}
					t = identifier(lexeme, lexeme, line, i)
				}
				res = append(res, t)
			} else {
				return nil, errors.New(fmt.Sprintf("Line %v Column %v: unexpected token %s", line, i, current()))
			}
		}
	}

	res = append(res, Token{tt: Eof})
	return res, nil
}
