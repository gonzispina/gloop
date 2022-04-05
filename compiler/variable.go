package compiler

type varType uint8

const (
	numberType varType = iota
	booleanType
)

func (vt varType) String() string {
	switch vt {
	case numberType:
		return "number"
	case booleanType:
		return "boolean"
	default:
		// Unreachable
		return ""
	}
}

func constantVarType(v interface{}) varType {
	if _, ok := v.(int); ok {
		return numberType
	} else {
		return booleanType
	}
}

type variable struct {
	name        string
	vt          varType
	initialized bool
	slot        byte
}
