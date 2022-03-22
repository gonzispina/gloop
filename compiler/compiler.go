package compiler

import (
	"github.com/gonzispina/gloop/vm"
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

func (p *Compiler) isAtEnd() bool {
	return p.counter == len(p.tokens) || p.tokens[p.counter].tt == Eof
}

func (p *Compiler) advance() Token {
	if p.isAtEnd() {
		return Token{tt: Eof}
	}
	p.counter++
	return p.tokens[p.counter-1]
}

func (p *Compiler) peek() Token {
	if p.isAtEnd() {
		return Token{tt: Eof}
	}
	return p.tokens[p.counter]
}

func (p *Compiler) match(tt tokenType) bool {
	t := p.peek()
	if t.tt == tt {
		p.counter++
		return true
	}

	return false
}

func (p *Compiler) expression() (interface{}, error) {
	return nil, nil
}

func (p *Compiler) procedureCall() error {
	return nil
}

func (p *Compiler) declareVariable(name string) *variable {
	slot := p.chunk.AddLocal()
	p.vars[name] = &variable{
		name:        name,
		vt:          0,
		initialized: false,
		slot:        slot,
	}
	return p.vars[name]
}

func (p *Compiler) varAssignment() error {
	t := p.peek()
	if t.tt == Eof {
		return unexpectedEndOfFileErr(p.peek())
	}

	name := t.value.(string)
	v, ok := p.vars[name]
	if !ok {
		v = p.declareVariable(name)
	}

	if !p.match(LeftArrow) {
		return expectedAssignmentOperatorErr(p.peek())
	}

	value, err := p.expression()
	if err != nil {
		return err
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Bool:
		if v.initialized && v.vt != boolean {
			return assignErr(p.peek(), v.vt, boolean)
		}
		v.initialized = true
		v.vt = boolean
	case reflect.Int64:
		if v.initialized && v.vt != number {
			return assignErr(p.peek(), v.vt, number)
		}
		v.initialized = true
		v.vt = number
	}

	p.chunk.Write(vm.OpSet.Byte(), t.line, v.slot)
	return nil
}

func (p *Compiler) statement() error {
	if p.match(If) {
		// return p.ifStatement()
	} else if p.match(Loop) {
		// return p.loopStatement()
	} else if p.match(Identifier) {
		return p.varAssignment()
	}

	return unexpectedTokenErr(p.peek())
}

func (p *Compiler) synchronize() {
	for {
		t := p.advance()
		if t.tt == Eof {
			break
		}

		if t.tt != End {
			continue
		}

		if p.match(Procedure) || p.match(Loop) || p.match(If) {
			break
		}
	}
}

func (p *Compiler) Parse() (vm.Chunk, []error) {
	var errs []error
	p.counter = 0
	for p.counter <= len(p.tokens) {
		err := p.statement()
		if err != nil {
			errs = append(errs, err)
			p.synchronize()
		}
	}

	if len(errs) != 0 {
		return vm.Chunk{}, errs
	}

	return p.chunk, nil
}
