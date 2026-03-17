package vm

import (
	"fmt"
	"math"
)

const defaultStackSize = 1024

type binaryOpHandler func(v1, v2 LValue) (LValue, error)

type Vm struct {
	stack     []LValue
	codes     []OpCode
	constants []LValue
	globals   map[string]LValue
	base      uint // stack base
	top       uint // stack top
	pc        uint // program counter

	// extra utils files
	binaryOpTable map[OpCode]binaryOpHandler
}

func NewVm() *Vm {
	vm := &Vm{
		stack:     make([]LValue, defaultStackSize),
		codes:     make([]OpCode, 0),
		constants: make([]LValue, 0),
		globals:   make(map[string]LValue),
		base:      0,
		top:       0,
		pc:        0,
	}
	vm.binaryOpTable = map[OpCode]binaryOpHandler{
		OpAdd: vm.opAdd,
		OpSub: vm.opSub,
		OpMul: vm.opMul,
		OpDiv: vm.opDiv,
		OpPow: vm.opPow,
	}
	return vm
}

func (v *Vm) WriteGlobals(key string, val LValue) {
	v.globals[key] = val
}

func (v *Vm) WriteConstant(val LValue) uint8 {
	v.constants = append(v.constants, val)
	return uint8(len(v.constants) - 1)
}

func (v *Vm) WriteCode(code OpCode, operands ...uint8) {
	v.codes = append(v.codes, code)
	for _, op := range operands {
		v.codes = append(v.codes, OpCode(op))
	}
}

func (v *Vm) Execute() error {
	for {
		if v.pc >= uint(len(v.codes)) {
			break
		}
		switch v.codes[v.pc] {
		case OpPrint: // only used for debugging
			val := v.stPop()
			fmt.Println(val.String())
			v.pc++
		case OpConstant:
			v.pc++
			idx := int(v.codes[v.pc])
			if idx >= len(v.constants) {
				return fmt.Errorf("invalid constant index")
			}
			v.stPush(v.constants[idx])
			v.pc++
		case OpAdd, OpSub, OpMul, OpDiv, OpPow:
			handler, ok := v.binaryOpTable[v.codes[v.pc]]
			if !ok {
				return fmt.Errorf("cannot find binary operator handler")
			}
			v1 := v.stPop()
			v2 := v.stPop()
			res, err := handler(v1, v2)
			if err != nil {
				return err
			}
			v.stPush(res)
			v.pc++
		case OpGetGlobal:
			v.pc++
			idx := int(v.codes[v.pc])
			if idx >= len(v.constants) {
				return fmt.Errorf("invalid constant index")
			}
			if v.constants[idx].Type() != LTString {
				return fmt.Errorf("invalid global name")
			}
			key := v.constants[idx].(*LString).Value
			val, ok := v.globals[key]
			if !ok {
				return fmt.Errorf("invalid global variable: %s", key)
			}
			v.stPush(val)
			v.pc++
		case OpSetGlobal:
			v.pc++
			idx := int(v.codes[v.pc])
			if idx >= len(v.constants) {
				return fmt.Errorf("invalid constant index")
			}
			if v.constants[idx].Type() != LTString {
				return fmt.Errorf("invalid global name")
			}
			key := v.constants[idx].(*LString).Value
			val := v.stPop()
			v.globals[key] = val
			v.pc++
		default:
			panic("invalid opcode")
		}
	}
	return nil
}

func (v *Vm) stPop() LValue {
	v.top--
	idx := v.base + v.top
	if int(idx) < 0 {
		panic("stack underflow")
	}
	val := v.stack[idx]
	v.stack[idx] = nil
	return val
}

func (v *Vm) stPush(val LValue) {
	idx := v.base + v.top
	if int(idx) >= len(v.stack) {
		stack := make([]LValue, 2*len(v.stack))
		copy(stack, v.stack)
		v.stack = stack
	}
	v.stack[idx] = val
	v.top++
}

func (v *Vm) convToLNumber(val LValue) (*LNumber, error) {
	switch val.Type() {
	case LTNumber:
		return val.(*LNumber), nil
	default:
		return nil, fmt.Errorf("cannot convert to number")
	}
}

func (v *Vm) opAdd(v1, v2 LValue) (LValue, error) {
	f1, err := v.convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := v.convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value + f2.Value,
	}, nil
}

func (v *Vm) opSub(v1, v2 LValue) (LValue, error) {
	f1, err := v.convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := v.convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value - f2.Value,
	}, nil
}

func (v *Vm) opMul(v1, v2 LValue) (LValue, error) {
	f1, err := v.convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := v.convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value * f2.Value,
	}, nil
}

func (v *Vm) opDiv(v1, v2 LValue) (LValue, error) {
	f1, err := v.convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := v.convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value / f2.Value,
	}, nil
}

func (v *Vm) opPow(v1, v2 LValue) (LValue, error) {
	f1, err := v.convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := v.convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: float32(math.Pow(float64(f1.Value), float64(f2.Value))),
	}, nil
}
