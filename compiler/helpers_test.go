package compiler_test

import (
	"github.com/gonzispina/gloop/compiler"
	"github.com/gonzispina/gloop/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func compile(t *testing.T, text string) (vm.Chunk, []error) {
	c := getCompiler(t, text)
	return c.Compile()
}

func getCompiler(t *testing.T, text string) *compiler.Compiler {
	tokens, err := compiler.Lexer(text)
	require.Nil(t, err)

	return compiler.New(tokens)
}

func assertErrContains(t *testing.T, errs []error, code compiler.ErrCode) {
	require.Equal(t, 1, len(errs))
	assert.True(t, strings.Contains(errs[0].Error(), string(code)))
}
