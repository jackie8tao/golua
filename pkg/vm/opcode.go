package vm

// OpCode
const (
	OpAdd       uint8 = iota + 128 // +
	OpSub                          // -
	OpMul                          // *
	OpDiv                          // /
	OpPow                          // ^
	OpConstant                     // load constant
	OpGetGlobal                    // get global variable
	OpSetGlobal                    // set global variable
	OpGetLocal                     // get local variable
	OpSetLocal                     // set local variable
	OpCall                         // function call
	OpReturn                       // return function
	OpPrint                        // for debugging
)

var opCodeNames = map[uint8]string{
	OpAdd:       "OpAdd",
	OpSub:       "OpSub",
	OpMul:       "OpMul",
	OpDiv:       "OpDiv",
	OpPow:       "OpPow",
	OpConstant:  "OpConstant",
	OpGetGlobal: "OpGetGlobal",
	OpSetGlobal: "OpSetGlobal",
	OpGetLocal:  "OpGetLocal",
	OpSetLocal:  "OpSetLocal",
	OpCall:      "OpCall",
	OpReturn:    "OpReturn",
	OpPrint:     "OpPrint",
}
