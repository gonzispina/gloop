package compiler

import (
	"github.com/gonzispina/gloop/vm"
	"path"
	"reflect"
)

func NewCompiler(tokens []Token) *Compiler {
	return &Compiler{
		chunk:   vm.NewChunk(),
		tokens:  tokens,
		counter: 0,
	}
}

type Compiler struct {
	chunk   vm.Chunk
	tokens  []Token
	counter int
	vars    map[string]*variable
}

func (c *Compiler) isAtEnd() bool {
	return c.counter == len(c.tokens) || c.tokens[c.counter].tt == Eof
}

func (c *Compiler) advance() Token {
	if c.isAtEnd() {
		return Token{tt: Eof}
	}
	c.counter++
	return c.tokens[c.counter-1]
}

func (c *Compiler) peek() Token {
	if c.isAtEnd() {
		return Token{tt: Eof}
	}
	return c.tokens[c.counter]
}

func (c *Compiler) match(tt tokenType) bool {
	t := c.peek()
	if t.tt == tt {
		c.counter++
		return true
	}

	return false
}

func (c *Compiler) expression() (interface{}, error) {
	if c.peek().tt == Identifier {

	}

	switch c.peek().tt {
	case Equal:
	case Greater:
	case Or:
	case

	}
}

func (c *Compiler) procedureCall() (interface{}, error) {
	return nil, nil
}

func (c *Compiler) declareVariable(name string) *variable {
	slot := c.chunk.AddLocal()
	c.vars[name] = &variable{
		name:        name,
		vt:          0,
		initialized: false,
		slot:        slot,
	}
	return c.vars[name]
}

func (c *Compiler) varAssignment() error {
	t := c.peek()
	if t.tt == Eof {
		return unexpectedEndOfFileErr(c.peek())
	}

	name := t.value.(string)
	v, ok := c.vars[name]
	if !ok {
		v = c.declareVariable(name)
	}

	if !c.match(LeftArrow) {
		return expectedAssignmentOperatorErr(c.peek())
	}

	value, err := c.expression()
	if err != nil {
		return err
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Bool:
		if v.initialized && v.vt != boolean {
			return assignErr(c.peek(), v.vt, boolean)
		}
		v.initialized = true
		v.vt = boolean
	case reflect.Int64:
		if v.initialized && v.vt != number {
			return assignErr(c.peek(), v.vt, number)
		}
		v.initialized = true
		v.vt = number
	}

	c.chunk.Write(vm.OpSet.Byte(), t.line, v.slot)
	return nil
}

func (c *Compiler) statement() error {
	if c.match(If) {
		// return c.ifStatement()
	} else if c.match(Loop) {
		// return c.loopStatement()
	} else if c.match(Identifier) {
		return c.varAssignment()
	}

	return unexpectedTokenErr(c.peek())
}

func (c *Compiler) synchronize() {
	for {
		t := c.advance()
		if t.tt == Eof {
			break
		}

		if t.tt != End {
			continue
		}

		if c.match(Procedure) || c.match(Loop) || c.match(If) {
			break
		}
	}
}

func (c *Compiler) Parse() (vm.Chunk, []error) {
	var errs []error
	c.counter = 0
	for c.counter <= len(c.tokens) {
		err := c.statement()
		if err != nil {
			errs = append(errs, err)
			c.synchronize()
		}
	}

	if len(errs) != 0 {
		return vm.Chunk{}, errs
	}

	return c.chunk, nil
}
