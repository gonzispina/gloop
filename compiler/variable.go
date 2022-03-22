package compiler

type varType uint8

const (
	number varType = iota
	boolean
)

func (vt varType) String() string {
	switch vt {
	case number:
		return "number"
	case boolean:
		return "boolean"
	default:
		// Unreachable
		return ""
	}
}

type variable struct {
	name        string
	vt          varType
	initialized bool
	slot        byte
}
