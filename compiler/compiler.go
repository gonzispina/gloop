package compiler

import (
	"github.com/gonzispina/gloop/vm"
	"time"
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

func (c *Compiler) parsePrecedence(previous Token, precedence Precedence) error {
	current := c.advance()
	prefixRule := getRule(c, previous.tt).prefix
	if prefixRule == nil {
		return expectedExpressionErr(previous)
	}

	err := prefixRule()
	if err != nil {
		return err
	}

	for rule := getRule(c, current.tt); rule.precedence > precedence; {
		current = c.advance()
		err := rule.infix()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) unary() error {
	t := c.peek()
	err := c.parsePrecedence(t, precedenceUnary)
	if err != nil {
		return err
	}

	c.chunk.Append(vm.OpNot.Byte(), t.line)
	return nil
}

func (c *Compiler) binary() error {
	t := c.peek()
	rule := getRule(c, t.tt)
	err := c.parsePrecedence(t, rule.precedence+1)
	if err != nil {
		return err
	}

	switch t.tt {
	case Plus:
		c.chunk.Append(vm.OpAdd.Byte(), t.line)
		break
	case Star:
		c.chunk.Append(vm.OpMultiply.Byte(), t.line)
		break
	case Equal:
		c.chunk.Append(vm.OpEqual.Byte(), t.line)
		break
	case Greater:
		c.chunk.Append(vm.OpGreater.Byte(), t.line)
		break
	case Lesser:
		c.chunk.Append(vm.OpLesser.Byte(), t.line)
		break
	case GreaterEqual:
		c.chunk.Append(vm.OpNot.Byte(), t.line)
		c.chunk.Append(vm.OpLesser.Byte(), t.line)
		break
	case LesserEqual:
		c.chunk.Append(vm.OpNot.Byte(), t.line)
		c.chunk.Append(vm.OpGreater.Byte(), t.line)
		break
	}

	return nil
}

func (c *Compiler) expression() error {
	return c.parsePrecedence(c.peek(), precedenceAssigment)
}

func (c *Compiler) grouping() error {
	t := c.advance()
	if err := c.expression(); err != nil {
		return err
	}

	t = c.advance()
	if t.tt != RightParen {
		return expectedRightParenthesisErr(t)
	}

	return nil
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
	t := c.advance()

	name := t.value.(string)
	v, ok := c.vars[name]
	if !ok {
		v = c.declareVariable(name)
	}

	if !c.match(LeftArrow) {
		return expectedAssignmentOperatorErr(c.peek())
	}

	err := c.expression()
	if err != nil {
		return err
	}

	/*
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
	*/

	c.chunk.Append(vm.OpSet.Byte(), t.line, v.slot)
	return nil
}

func (c *Compiler) ifStatement() error {
	if err := c.expression(); err != nil {
		return err
	}

	if !c.match(Then) {
		return expectedThenErr(c.peek())
	}

	thenJumpOffset := c.chunk.EmitJump(vm.OpJumpIf, c.peek().line)
	if err := c.statement(); err != nil {
		return err
	}

	if c.match(Else) {

	}

	if err := c.chunk.PatchJump(thenJumpOffset); err != nil {
		return blockIsTooLargeErr(c.peek())
	}

	for c.peek().tt != End

}

func (c *Compiler) statement() error {
	if c.match(If) {
		return c.ifStatement()
	} else if c.match(Loop) {
		// return c.loopStatement()
	} else if c.peek().tt == Identifier {
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
