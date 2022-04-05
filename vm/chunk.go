package vm

import "errors"

func NewChunk() Chunk {
	return Chunk{
		ip:           0,
		instructions: []byte{},
		locals:       map[byte]byte{},
	}
}

// Chunk for instructions
type Chunk struct {
	ip           uint32
	instructions []byte
	line         []int
	locals       map[byte]byte
	localCount   int
}

func (c Chunk) InstructionsCount() int {
	return len(c.instructions)
}

func (c Chunk) AddLocal() byte {
	c.localCount++
	c.locals[byte(c.localCount)] = 0x0
	return byte(c.localCount)
}

func (c Chunk) Write(index int, b byte) error {
	if index >= len(c.instructions) {
		panic("invalid chunk index")
	}

	c.instructions[index] = b
	return nil
}

func (c Chunk) EmitJump(code OpCode, line int) int {
	return c.Append(code.Byte(), line, 0xff, 0xff)
}

func (c Chunk) PatchJump(offset int, to ...int) error {
	jump := len(c.instructions) - 2 - offset
	if jump > 256 {
		return errors.New("block is too large")
	}

	if len(to) > 0 {
		jump = to[0]
	}

	c.instructions[offset] = byte(jump >> 8 & 0xff)
	c.instructions[offset+1] = byte(jump & 0xff)

	return nil
}

func (c Chunk) Append(b byte, line int, more ...byte) int {
	c.instructions = append(c.instructions, b)
	c.line = append(c.line, line)

	if len(more) > 0 {
		for _, b := range more {
			c.instructions = append(c.instructions, b)
			c.line = append(c.line, line)
		}
	}

	return len(c.instructions) - 1
}

func (c Chunk) Read() byte {
	return c.instructions[c.ip]
}
