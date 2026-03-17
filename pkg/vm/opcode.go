package vm

type OpCode uint8

const (
	OpAdd       OpCode = iota + 128 // +
	OpSub                           // -
	OpMul                           // *
	OpDiv                           // /
	OpPow                           // ^
	OpConstant                      // load constant
	OpGetGlobal                     // get global value
	OpSetGlobal                     // set global value
	OpPrint                         // for debugging
)
