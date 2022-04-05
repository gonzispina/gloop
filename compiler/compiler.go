package compiler

import (
	"github.com/gonzispina/gloop/vm"
)

func New(tokens []Token) *Compiler {
	return &Compiler{
		chunk:   vm.NewChunk(),
		tokens:  tokens,
		counter: 0,
		vars:    map[string]*variable{},
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

func (c *Compiler) parsePrecedence(previous Token, precedence Precedence) (interface{}, error) {
	prefixRule := getRule(c, previous.tt).prefix
	if prefixRule == nil {
		return nil, expectedExpressionErr(previous)
	}

	v, err := prefixRule()
	if err != nil {
		return nil, err
	}

	current := c.advance()
	for rule := getRule(c, current.tt); rule.precedence > precedence; {
		current = c.advance()
		_, err := rule.infix()
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (c *Compiler) constant() (interface{}, error) {
	t := c.peek()
	var res varType
	if v, ok := t.value.(int64); ok {
		c.chunk.Append(vm.OpPush.Byte(), t.line, byte(int64(v)))
		res = numberType
	} else if b, ok := t.value.(bool); ok {
		if b {
			c.chunk.Append(vm.OpPush.Byte(), t.line, byte(1))
		} else {
			c.chunk.Append(vm.OpPush.Byte(), t.line, byte(0))
		}
		res = booleanType
	}

	return res, nil
}

func (c *Compiler) unary() (interface{}, error) {
	t := c.peek()
	v, err := c.parsePrecedence(t, precedenceUnary)
	if err != nil {
		return nil, err
	}

	c.chunk.Append(vm.OpNot.Byte(), t.line)
	return v, nil
}

func (c *Compiler) binary() (interface{}, error) {
	t := c.peek()
	rule := getRule(c, t.tt)
	v, err := c.parsePrecedence(t, rule.precedence+1)
	if err != nil {
		return nil, err
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

	return v, nil
}

func (c *Compiler) expression() (interface{}, error) {
	return c.parsePrecedence(c.peek(), precedenceAssigment)
}

func (c *Compiler) grouping() (interface{}, error) {
	t := c.advance()
	v, err := c.expression()
	if err != nil {
		return nil, err
	}

	t = c.advance()
	if t.tt != RightParen {
		return nil, expectedRightParenthesisErr(t)
	}

	return v, nil
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

func (c *Compiler) varEvaluation() (interface{}, error) {
	t := c.advance()

	name := t.value.(string)
	_, ok := c.vars[name]
	if !ok {
		return nil, undefinedVariableErr(t, name)
	}

	return nil, nil
}

func (c *Compiler) varAssignment() (interface{}, error) {
	t := c.advance()

	name := t.value.(string)
	v, ok := c.vars[name]
	if !ok {
		v = c.declareVariable(name)
	}

	if !c.match(LeftArrow) {
		return nil, expectedAssignmentOperatorErr(c.peek())
	}

	e, err := c.expression()
	if err != nil {
		return nil, err
	}

	switch e.(varType) {
	case booleanType:
		if v.initialized && v.vt != booleanType {
			return nil, invalidTypeErr(c.peek(), v.vt, booleanType)
		}
		v.initialized = true
		v.vt = booleanType
	case numberType:
		if v.initialized && v.vt != numberType {
			return nil, invalidTypeErr(c.peek(), v.vt, numberType)
		}
		v.initialized = true
		v.vt = numberType
	}

	c.chunk.Append(vm.OpSet.Byte(), t.line, v.slot)
	return nil, nil
}

func (c *Compiler) ifStatement() (interface{}, error) {
	t := c.peek()

	v, err := c.expression()
	if err != nil {
		return nil, err
	}

	if expressionType := v.(varType); expressionType != booleanType {
		return nil, booleanExpressionNeededErr(t)
	}

	if !c.match(Then) {
		return nil, expectedThenErr(c.peek())
	}

	thenJumpOffset := c.chunk.EmitJump(vm.OpJumpIfFalse, c.peek().line)
	for c.peek().tt != EndIf && c.peek().tt != Else {
		if _, err := c.statement(); err != nil {
			return nil, err
		}
	}

	for c.match(Else) {
		if c.peek().tt == If {
			t = c.advance()
			previousOffset := thenJumpOffset
			if v, err = c.expression(); err != nil {
				return nil, err
			}

			if expressionType := v.(varType); expressionType != booleanType {
				return nil, booleanExpressionNeededErr(t)
			}

			if !c.match(Then) {
				return nil, expectedThenErr(c.peek())
			}

			thenJumpOffset = c.chunk.EmitJump(vm.OpJumpIfFalse, c.peek().line)
			if _, err := c.statement(); err != nil {
				return nil, err
			}

			c.chunk.Append(vm.OpPop.Byte(), c.peek().line)
			if err := c.chunk.PatchJump(previousOffset); err != nil {
				return nil, blockIsTooLargeErr(c.peek())
			}

			for c.peek().tt != EndIf && c.peek().tt != Else {
				if _, err := c.statement(); err != nil {
					return nil, err
				}
			}

			continue
		}

		c.chunk.Append(vm.OpPop.Byte(), c.peek().line)
		if err := c.chunk.PatchJump(thenJumpOffset); err != nil {
			return nil, blockIsTooLargeErr(c.peek())
		}

		for c.peek().tt != EndIf {
			if _, err := c.statement(); err != nil {
				return nil, err
			}
		}
	}

	if !c.match(EndIf) {
		return nil, expectedEndIfErr(c.peek())
	}

	return nil, nil
}

func (c *Compiler) loopStatement() (interface{}, error) {
	t := c.peek()

	loopStart := c.chunk.InstructionsCount()

	v, err := c.expression()
	if err != nil {
		return nil, err
	}

	if expressionType := v.(varType); expressionType != numberType {
		return nil, numberExpressionNeededErr(t)
	}

	exitJump := c.chunk.EmitJump(vm.OpJumpIfFalse, c.peek().line)
	if !c.match(Times) {
		return nil, expectedTimesErr(c.peek())
	}

	var abortsToPatch []int
	for c.peek().tt != EndLoop {
		if c.match(AbortLoop) {
			abortsToPatch = append(
				abortsToPatch,
				c.chunk.EmitJump(vm.OpJump, c.peek().line),
			)
		} else {
			if _, err := c.statement(); err != nil {
				return nil, err
			}
		}
	}

	if !c.match(EndLoop) {
		return nil, expectedEndLoopErr(c.peek())
	}

	loopBack := c.chunk.EmitJump(vm.OpJump, c.peek().line)
	if err := c.chunk.PatchJump(loopBack, loopStart); err != nil {
		return nil, blockIsTooLargeErr(c.peek())
	}

	c.chunk.Append(vm.OpPop.Byte(), c.peek().line)
	if err := c.chunk.PatchJump(exitJump); err != nil {
		return nil, blockIsTooLargeErr(c.peek())
	}

	for _, jump := range abortsToPatch {
		if err := c.chunk.PatchJump(jump); err != nil {
			return nil, blockIsTooLargeErr(c.peek())
		}
	}

	return nil, nil
}

func (c *Compiler) statement() (interface{}, error) {
	if c.match(If) {
		return c.ifStatement()
	} else if c.match(Loop) {
		return c.loopStatement()
	} else if c.peek().tt == Identifier {
		return c.varAssignment()
	}

	return nil, unexpectedTokenErr(c.peek())
}

func (c *Compiler) synchronize() {
	for {
		t := c.advance()
		if t.tt == Eof {
			break
		}

		if c.match(EndProcedure) || c.match(AbortLoop) || c.match(EndIf) {
			break
		}
	}
}

func (c *Compiler) Compile() (vm.Chunk, []error) {
	var errs []error
	c.counter = 0
	for c.counter < len(c.tokens)-1 {
		_, err := c.statement()
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
