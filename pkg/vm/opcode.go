package vm

type OpCode uint8

const (
	OpAdd OpCode = iota + 128 // +
	OpSub                     // -
	OpMul                     // *
	OpDiv                     // /
	OpPow                     // ^
	OpConstant
	OpPrint // for debugging
)
