package parser

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
			{tt: Define},
			{tt: Procedure},
			{tt: Literal, value: "MINUS"},
			{tt: LeftSquareBracket},
			{tt: Literal, value: "M"},
			{tt: Comma},
			{tt: Literal, value: "N"},
			{tt: RightSquareBracket},
			{tt: If},
			{tt: Literal, value: "M"},
			{tt: Lesser},
			{tt: Literal, value: "N"},
			{tt: Then},
			{tt: Quit},
			{tt: Procedure},
			{tt: End},
			{tt: If},
			{tt: Loop},
			{tt: Literal, value: "M"},
			{tt: Plus},
			{tt: Literal, value: int64(1)},
			{tt: Times},
			{tt: If},
			{tt: Output},
			{tt: Plus},
			{tt: Literal, value: "N"},
			{tt: Equal},
			{tt: Literal, value: "M"},
			{tt: Then},
			{tt: Abort},
			{tt: Loop},
			{tt: End},
			{tt: If},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Output},
			{tt: Plus},
			{tt: Literal, value: int64(1)},
			{tt: End},
			{tt: Loop},
			{tt: End},
			{tt: Procedure},
		}

		res, err := Lexer(text)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
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
			{tt: Define},
			{tt: Procedure},
			{tt: Literal, value: "ISEVEN?"},
			{tt: LeftSquareBracket},
			{tt: Literal, value: "N"},
			{tt: RightSquareBracket},

			{tt: If},
			{tt: Literal, value: "N"},
			{tt: Lesser},
			{tt: Literal, value: int64(2)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Literal, value: true},
			{tt: Quit},
			{tt: Procedure},
			{tt: End},
			{tt: If},

			{tt: Loop},
			{tt: Literal, value: "N"},
			{tt: Times},

			{tt: Literal, value: "N"},
			{tt: LeftArrow},
			{tt: Literal, value: "MINUS"},
			{tt: LeftSquareBracket},
			{tt: Literal, value: "N"},
			{tt: Comma},
			{tt: Literal, value: int64(2)},
			{tt: RightSquareBracket},

			{tt: If},
			{tt: Literal, value: "N"},
			{tt: Equal},
			{tt: Literal, value: int64(1)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Literal, value: false},
			{tt: Quit},
			{tt: Procedure},

			{tt: Else},
			{tt: If},
			{tt: Literal, value: "N"},
			{tt: Equal},
			{tt: Literal, value: int64(0)},
			{tt: Then},
			{tt: Output},
			{tt: LeftArrow},
			{tt: Literal, value: true},
			{tt: Quit},
			{tt: Procedure},
			{tt: End},
			{tt: If},

			{tt: End},
			{tt: Loop},
			{tt: End},
			{tt: Procedure},
		}

		res, err := Lexer(text)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})
}
