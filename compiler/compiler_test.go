package compiler_test

import (
	"github.com/gonzispina/gloop/compiler"
	"testing"
)

func TestCompiler_Compile_Assignments(t *testing.T) {
	t.Run("Bad assignment operator return an expected assignment operator err", func(t *testing.T) {
		text := `
			OUTPUT = YES
		`

		_, errs := compile(t, text)
		assertErrContains(t, errs, compiler.ExpectedAssignmentOperatorErrCode)
	})

	t.Run("Bad assignment type returns an invalid type error", func(t *testing.T) {
		text := `
			OUTPUT <- YES
			OUTPUT <- 1
		`

		_, errs := compile(t, text)
		assertErrContains(t, errs, compiler.InvalidTypeErrCode)

		text = `
			OUTPUT <- 1
			OUTPUT <- NO
		`

		_, errs = compile(t, text)
		assertErrContains(t, errs, compiler.InvalidTypeErrCode)
	})

	t.Run("Uninitialized variable usage returns uninitialized variable error", func(t *testing.T) {
		text := `
			OUTPUT <- N
		`

		_, errs := compile(t, text)
		assertErrContains(t, errs, compiler.UndefinedVariableErrCode)
	})
	
	t.Run("Initialization of a variable with itself returns an uninitialized variable error", func(t *testing.T) {
		text := `
			N <- N + 1 
			OUTPUT <- N
		`

		_, errs := compile(t, text)
		assertErrContains(t, errs, compiler.UndefinedVariableErrCode)
	})
}

/*
func TestCompiler_Compile_If_Statements(t *testing.T) {
	t.Run("If statements", func(t *testing.T) {
		t.Run("Unterminated 'If - Then' returns unexpected end of file err", func(t *testing.T) {
			text := `
				N <- 1
				IF N < 2 THEN
					OUTPUT <- YES
			`
		})

		t.Run("Unterminated 'Else If - Then' returns unexpected end of file err", func(t *testing.T) {
			text := `
				N <- 1
				IF N < 2 THEN
					OUTPUT <- YES
				ELSE IF N < 3 THEN
					OUTPUT <- NO
			`
		})

		t.Run("Unterminated 'Else' returns unexpected end of file err", func(t *testing.T) {
			text := `
				N <- 1
				IF N < 2 THEN
					OUTPUT <- YES
				ELSE
					OUTPUT <- NO
			`
		})

		t.Run("Evaluating non boolean expression return err", func(t *testing.T) {
			text := `
				N <- 1
				IF N THEN
					OUTPUT <- YES
				END IF
			`
		})

		t.Run("No else 'OK' scenario", func(t *testing.T) {
			t.Run("Jumps if the condition is false", func(t *testing.T) {
				text := `
					N <- 1
					IF N < 2 THEN
						N <- N + 1
					END IF
					OUTPUT <- N
				`
			})

			t.Run("Doesn't jump if the condition is true", func(t *testing.T) {
				text := `
					N <- 1
					IF N < 2 THEN
						N <- N + 1
					END IF
					OUTPUT <- N
				`
			})
		})

		t.Run("Else case scenario", func(t *testing.T) {
			t.Run("Jumps to the else branch if the condition is false", func(t *testing.T) {
				text := `
					N <- 1
					IF N < 2 THEN
						N <- N + 1
					ELSE
						N <- N - 1
					END IF
					OUTPUT <- N
				`
			})

			t.Run("Runs the main then branch if the condition is true", func(t *testing.T) {
				text := `
					N <- 1
					IF N < 2 THEN
						N <- N + 1
					ELSE
						N <- N - 1
					END IF
					OUTPUT <- N
				`
			})
		})

	})
}
*/
