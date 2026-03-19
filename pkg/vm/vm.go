package vm

import (
	"fmt"
	"math"
)

const defaultStackSize = 1024

type binaryOpHandler func(v1, v2 LValue) (LValue, error)

type Vm struct {
	stack         []LValue
	base          uint // stack base
	top           uint // stack top
	fn            *FuncProto
	binaryOpTable map[uint8]binaryOpHandler
}

func NewVm(funcProto *FuncProto) *Vm {
	vm := &Vm{
		stack: make([]LValue, defaultStackSize),
		base:  0,
		top:   uint(funcProto.maxLocals),
		fn:    funcProto,
	}
	vm.binaryOpTable = map[uint8]binaryOpHandler{
		OpAdd: vm.opAdd,
		OpSub: vm.opSub,
		OpMul: vm.opMul,
		OpDiv: vm.opDiv,
		OpPow: vm.opPow,
	}
	return vm
}

func (v *Vm) Execute() error {
	for {
		if v.fn.pc >= len(v.fn.codes) {
			break
		}
		switch v.fn.codes[v.fn.pc] {
		case OpPrint: // only used for debugging
			val := v.stPop()
			fmt.Println(val.String())
			v.fn.pc++
		case OpConstant:
			args := v.fnOpArg(2)
			idx := convToUint16([2]uint8{args[0], args[1]})
			val, err := v.getFnConstant(idx)
			if err != nil {
				return err
			}
			v.stPush(val)
			v.fn.pc++
		case OpAdd, OpSub, OpMul, OpDiv, OpPow:
			handler, ok := v.binaryOpTable[v.fn.codes[v.fn.pc]]
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
			v.fn.pc++
		case OpGetGlobal:
			args := v.fnOpArg(2)
			idx := convToUint16([2]uint8{args[0], args[1]})
			key, err := v.getFnConstant(idx)
			if err != nil {
				return err
			}
			if key.Type() != LTString {
				return fmt.Errorf("invalid global name type")
			}
			val, ok := v.fn.globals[key.(*LString).Value]
			if !ok {
				return fmt.Errorf("cannot find global value")
			}
			v.stPush(val)
			v.fn.pc++
		case OpSetGlobal:
			args := v.fnOpArg(2)
			idx := convToUint16([2]uint8{args[0], args[1]})
			key, err := v.getFnConstant(idx)
			if err != nil {
				return err
			}
			if key.Type() != LTString {
				return fmt.Errorf("invalid global name type")
			}
			val := v.stPop()
			v.fn.globals[key.(*LString).Value] = val
			v.fn.pc++
		case OpSetLocal:
			args := v.fnOpArg(2)
			idx := convToUint16([2]uint8{args[0], args[1]})
			val := v.stPop()
			v.stack[v.base+uint(idx)] = val
			v.fn.pc++
		case OpGetLocal:
			args := v.fnOpArg(2)
			idx := convToUint16([2]uint8{args[0], args[1]})
			val := v.stack[v.base+uint(idx)]
			v.stPush(val)
			v.fn.pc++
		default:
			panic("invalid opcode")
		}
	}
	return nil
}

func (v *Vm) fnOpArg(n int) []uint8 {
	res := make([]uint8, 0)
	for i := 0; i < n; i++ {
		v.fn.pc++
		res = append(res, v.fn.codes[v.fn.pc])
	}
	return res
}

// used for vm to get constant by index
func (v *Vm) getFnConstant(idx uint16) (LValue, error) {
	if int(idx) >= len(v.fn.constants) {
		return nil, fmt.Errorf("constant index overflow")
	}
	return v.fn.constants[idx], nil
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
		Value: math.Pow(f1.Value, f2.Value),
	}, nil
}
