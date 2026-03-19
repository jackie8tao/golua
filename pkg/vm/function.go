package vm

type FuncProto struct {
	codes     []uint8
	constants []LValue
	globals   map[string]LValue
	pc        int
	maxLocals int
}

func NewFuncProto() *FuncProto {
	return &FuncProto{
		codes:     make([]uint8, 0),
		constants: make([]LValue, 0),
		globals:   make(map[string]LValue),
		pc:        0,
		maxLocals: 0,
	}
}
