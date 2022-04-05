package compiler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("It returns the correct tokens for number procedure", func(t *testing.T) {
		text := `
			DEFINE PROCEDURE "MINUS" [M, N]
				IF M < N THEN
					QUIT PROCEDURE
				END IF
				LOOP M + 1 TIMES
					IF OUTPUT + N = M THEN
						ABORT LOOP
					END IF
					OUTPUT <- OUTPUT + 1
				END LOOP
			END PROCEDURE
		`

		expected := []Token{
			{tt: DefineProcedure},
			{tt: Identifier, value: "MINUS"},
			{tt: LeftSquareBracket},
			{tt: Identifier, value: "M"},
			{tt: Comma},
			{tt: Identifier, value: "N"},
			{tt: RightSquareBracket},
			{tt: If},
			{tt: Identifier, value: "M"},
			{tt: Lesser},
			{tt: Identifier, value: "N"},
			{tt: Then},
			{tt: QuitProcedure},
			{tt: EndIf},
			{tt: Loop},
			{tt: Identifier, value: "M"},
			{tt: Plus},
			{tt: Constant, value: int64(1)},
			{tt: Times},
			{tt: If},
			{tt: Output},
			{tt: Plus},
			{tt: Identifier, value: "N"},
			{tt: Equal},
			{tt: Identifier, value: "M"},
			{tt: Then},
			{tt: AbortLoop},
			{tt: EndIf},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Output},
			{tt: Plus},
			{tt: Constant, value: int64(1)},
			{tt: EndLoop},
			{tt: EndProcedure},
			{tt: Eof},
		}

		res, err := Lexer(text)
		assert.Nil(t, err)
		for i, tkn := range res {
			assert.Equal(t, expected[i].tt, tkn.tt)
		}
	})

	t.Run("It returns the correct tokens for test procedure", func(t *testing.T) {
		text := `
			DEFINE PROCEDURE "ISEVEN?" [N]
				IF N < 2 THEN
					OUTPUT <- YES
					QUIT PROCEDURE
				END IF

				LOOP N TIMES 
					N <- MINUS[N, 2]
					IF N = 1 THEN
						OUTPUT <- NO
						QUIT PROCEDURE
					ELSE IF N = 0 THEN
						OUTPUT <- YES
						QUIT PROCEDURE
					END IF
				END LOOP
			END PROCEDURE
		`

		expected := []Token{
			{tt: DefineProcedure},
			{tt: Identifier, value: "ISEVEN?"},
			{tt: LeftSquareBracket},
			{tt: Identifier, value: "N"},
			{tt: RightSquareBracket},

			{tt: If},
			{tt: Identifier, value: "N"},
			{tt: Lesser},
			{tt: Constant, value: int64(2)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Constant, value: true},
			{tt: QuitProcedure},
			{tt: EndIf},

			{tt: Loop},
			{tt: Identifier, value: "N"},
			{tt: Times},

			{tt: Identifier, value: "N"},
			{tt: LeftArrow},
			{tt: Identifier, value: "MINUS"},
			{tt: LeftSquareBracket},
			{tt: Identifier, value: "N"},
			{tt: Comma},
			{tt: Constant, value: int64(2)},
			{tt: RightSquareBracket},

			{tt: If},
			{tt: Identifier, value: "N"},
			{tt: Equal},
			{tt: Constant, value: int64(1)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Constant, value: false},
			{tt: QuitProcedure},

			{tt: Else},
			{tt: If},
			{tt: Identifier, value: "N"},
			{tt: Equal},
			{tt: Constant, value: int64(0)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Constant, value: true},
			{tt: QuitProcedure},
			{tt: EndIf},

			{tt: EndLoop},
			{tt: EndProcedure},
			{tt: Eof},
		}

		res, err := Lexer(text)
		assert.Nil(t, err)
		for i, tkn := range res {
			assert.Equal(t, expected[i].tt, tkn.tt)
		}
	})
}
