package vm

// OpCode
const (
	OpAdd       uint8 = iota + 128 // +, stack: (2) -> (1)
	OpSub                          // -, stack: (2) -> (1)
	OpMul                          // *, stack: (2) -> (1)
	OpDiv                          // /, stack: (2) -> (1)
	OpPow                          // ^, stack: (2) -> (1)
	OpConstant                     // load constant, stack: (0) -> 1
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
