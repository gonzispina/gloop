package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("It returns the correct tokens", func(t *testing.T) {
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
		for i, v := range res {
			if expected[i] != v {
				break
			}
			assert.Equal(t, expected[i], v)
		}
	})
}
