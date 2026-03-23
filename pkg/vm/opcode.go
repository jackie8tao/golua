package vm

// OpCode
const (
	OpAdd       uint8 = iota + 128 // +, a1 = pop(), a2 = pop(), push(a1 + a2)
	OpSub                          // -, a1 = pop(), a2 = pop(), push(a2 - a1)
	OpMul                          // *, a1 = pop(), a2 = pop(), push(a1 * a2)
	OpDiv                          // /, a1 = pop(), a2 = pop(), push(a2 / a1)
	OpPow                          // ^, a1 = pop(), a2 = pop(), push(a2 ^ a1)
	OpConstant                     // idx = opcode[pc], push(constants[idx])
	OpGetGlobal                    // idx = opcode[pc], push(globals[constants[idx]])
	OpSetGlobal                    // idx = opcode[pc], val = pop(), globals[constants[idx]] = val
	OpGetLocal                     // idx = opcode[pc], push(stack[bp+idx])
	OpSetLocal                     // idx = opcode[pc], val = pop(), stack[bp+idx] = val
	OpCall                         // argc = opcode[pc], nres=opcode[pc++],fn = stack[bp+sp-argc-1], push(fn(v1,v2,))
	OpReturn                       // return function
	OpLoadNil                      // push(nil)
	OpLoadTrue                     // push(true)
	OpLoadFalse                    // push(false)
	OpNewTable                     // push({})
	OpSetTable                     // idx = opcode[pc], val = pop(), table[constants[idx]] = val
	OpGetTable                     // idx = opcode[pc], push(table[constants[idx]])
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
	OpLoadNil:   "OpLoadNil",
	OpLoadTrue:  "OpLoadTrue",
	OpLoadFalse: "OpLoadFalse",
	OpNewTable:  "OpNewTable",
	OpSetTable:  "OpSetTable",
	OpGetTable:  "OpGetTable",
}
