package vm

type OpCode byte

func (o OpCode) Byte() byte {
	return byte(o)
}

const (
	OpAdd OpCode = iota
	OpMultiply
	OpEqual
	OpGreater
	OpLesser

	OpPush
	OpPop
	OpJump
	OpJumpIf

	OpSet
	OpGet
)
