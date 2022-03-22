package vm

func NewChunk() Chunk {
	return Chunk{ip: 0, instructions: []byte{}}
}

// Chunk for instructions
type Chunk struct {
	ip           uint32
	instructions []byte
	line         []int
	locals       map[byte]byte
	localCount   int
}

func (c Chunk) AddLocal() byte {
	c.localCount++
	c.locals[byte(c.localCount)] = 0x0
	return byte(c.localCount)
}

func (c Chunk) Write(b byte, line int, more ...byte) {
	c.instructions = append(c.instructions, b)
	c.line = append(c.line, line)

	if len(more) > 0 {
		for _, b := range more {
			c.instructions = append(c.instructions, b)
			c.line = append(c.line, line)
		}
	}
}

func (c Chunk) Read() byte {
	return c.instructions[c.ip]
}
