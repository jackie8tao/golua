package vm

type OpCode uint8

const (
	OpAdd       OpCode = iota + 128 // +
	OpSub                           // -
	OpMul                           // *
	OpDiv                           // /
	OpPow                           // ^
	OpConstant                      // load constant
	OpGetGlobal                     // get global variable
	OpSetGlobal                     // set global variable
	OpGetLocal                      // get local variable
	OpSetLocal                      // set local variable
	OpPrint                         // for debugging
)
