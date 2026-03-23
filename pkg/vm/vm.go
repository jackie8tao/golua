package vm

import (
	"fmt"
	"math"
)

const defaultStackSize = 1024

type binaryOpHandler func(v1, v2 LValue) (LValue, error)

type Vm struct {
	stack    []LValue
	bp       uint // stack bp
	sp       uint // stack sp
	funcs    map[string]LGFunction
	curFrame *CallFrame
}

var opTable = map[uint8]func(v *Vm) error{
	OpAdd:       opArithmetic,
	OpSub:       opArithmetic,
	OpMul:       opArithmetic,
	OpDiv:       opArithmetic,
	OpPow:       opArithmetic,
	OpSetLocal:  opSetLocal,
	OpGetLocal:  opGetLocal,
	OpConstant:  opConstant,
	OpSetGlobal: opSetGlobal,
	OpGetGlobal: opGetGlobal,
	OpCall:      opCall,
	OpReturn:    opReturn,
}

func NewVm(funcProto *FuncProto) *Vm {
	vm := &Vm{
		stack: make([]LValue, defaultStackSize),
		bp:    0,
		sp:    uint(funcProto.maxLocals),
		funcs: make(map[string]LGFunction),
		curFrame: &CallFrame{
			fn:     funcProto,
			parent: nil,
			bp:     0,
			sp:     0,
			pc:     0,
		},
	}
	return vm
}

func (v *Vm) RegisterModule(module string, funcs map[string]LGFunction) {
	for name, fn := range funcs {
		v.funcs[name] = fn
	}
}

func (v *Vm) StPop() LValue {
	v.sp--
	idx := v.bp + v.sp
	if int(idx) < 0 {
		panic("stack underflow")
	}
	val := v.stack[idx]
	v.stack[idx] = nil
	return val
}

func (v *Vm) StPush(val LValue) {
	idx := v.bp + v.sp
	if int(idx) >= len(v.stack) {
		stack := make([]LValue, 2*len(v.stack))
		copy(stack, v.stack)
		v.stack = stack
	}
	v.stack[idx] = val
	v.sp++
}

func (v *Vm) Execute() error {
	for {
		if v.curFrame.fn.pc >= len(v.curFrame.fn.codes) {
			break
		}
		handler, ok := opTable[v.curFrame.fn.codes[v.curFrame.fn.pc]]
		if !ok {
			return fmt.Errorf("not support opcode")
		}
		err := handler(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Vm) fnOpIdxArg() uint16 {
	args := v.fnOpArg(2)
	idx := convToUint16([2]uint8{args[0], args[1]})
	return idx
}

func (v *Vm) fnOpArg(n int) []uint8 {
	res := make([]uint8, 0)
	for i := 0; i < n; i++ {
		v.curFrame.fn.pc++
		res = append(res, v.curFrame.fn.codes[v.curFrame.fn.pc])
	}
	return res
}

// used for vm to get constant by index
func (v *Vm) getFnConstant(idx uint16) (LValue, error) {
	if int(idx) >= len(v.curFrame.fn.constants) {
		return nil, fmt.Errorf("constant index overflow")
	}
	return v.curFrame.fn.constants[idx], nil
}

func opConstant(v *Vm) error {
	idx := v.fnOpIdxArg()
	val, err := v.getFnConstant(idx)
	if err != nil {
		return err
	}
	v.StPush(val)
	v.curFrame.fn.pc++
	return nil
}

func opArithmetic(v *Vm) error {
	var binaryOpTable = map[uint8]binaryOpHandler{
		OpAdd: opAdd,
		OpSub: opSub,
		OpMul: opMul,
		OpDiv: opDiv,
		OpPow: opPow,
	}
	handler, ok := binaryOpTable[v.curFrame.fn.codes[v.curFrame.fn.pc]]
	if !ok {
		return fmt.Errorf("cannot find binary operator handler")
	}
	v1 := v.StPop()
	v2 := v.StPop()
	res, err := handler(v1, v2)
	if err != nil {
		return err
	}
	v.StPush(res)
	v.curFrame.fn.pc++
	return nil
}

func opGetGlobal(v *Vm) error {
	idx := v.fnOpIdxArg()
	key, err := v.getFnConstant(idx)
	if err != nil {
		return err
	}
	if key.Type() != LTString {
		return fmt.Errorf("invalid global name type")
	}
	lKey := key.(*LString)
	val, ok := v.curFrame.fn.globals[lKey.Value]
	if !ok {
		val, ok = v.funcs[lKey.Value]
		if !ok {
			return fmt.Errorf("cannot find global value: %s", lKey.Value)
		}
	}
	v.StPush(val)
	v.curFrame.fn.pc++
	return nil
}

func opSetGlobal(v *Vm) error {
	idx := v.fnOpIdxArg()
	key, err := v.getFnConstant(idx)
	if err != nil {
		return err
	}
	if key.Type() != LTString {
		return fmt.Errorf("invalid global name type")
	}
	val := v.StPop()
	v.curFrame.fn.globals[key.(*LString).Value] = val
	v.curFrame.fn.pc++
	return nil
}

func opSetLocal(v *Vm) error {
	idx := v.fnOpIdxArg()
	val := v.StPop()
	v.stack[v.bp+uint(idx)] = val
	v.curFrame.fn.pc++
	return nil
}

func opGetLocal(v *Vm) error {
	idx := v.fnOpIdxArg()
	val := v.stack[v.bp+uint(idx)]
	v.StPush(val)
	v.curFrame.fn.pc++
	return nil
}

func opCall(v *Vm) error {
	argCount := v.fnOpIdxArg()
	fn := v.stack[(v.bp+v.sp)-uint(argCount)-1]
	switch callee := fn.(type) {
	case *LFunction:
		newBp := v.sp - uint(argCount)
		v.curFrame.bp = v.bp
		newCallFrame := &CallFrame{
			fn:     callee.Fn,
			parent: v.curFrame,
			bp:     newBp,
			pc:     0,
		}
		v.bp = newBp
		v.sp = newBp + uint(newCallFrame.fn.maxLocals)
		v.curFrame = newCallFrame
	case LGFunction:
		retCount := uint(callee(v, int(argCount)))
		baseIdx := (v.bp + v.sp) - retCount - 1
		for i := uint(0); i < retCount; i++ {
			v.stack[baseIdx+i] = v.stack[baseIdx+i+1]
		}
		v.sp--
		v.curFrame.fn.pc++
	default:
		return fmt.Errorf("invalid function type")
	}
	return nil
}

func opReturn(v *Vm) error {
	retVal := v.StPop()
	v.sp = v.curFrame.bp - 1
	v.curFrame = v.curFrame.parent
	if v.curFrame == nil {
		return nil
	}
	v.bp = v.curFrame.bp
	v.StPush(retVal)
	v.curFrame.fn.pc++
	return nil
}

func opAdd(v1, v2 LValue) (LValue, error) {
	f1, err := convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value + f2.Value,
	}, nil
}

func opSub(v1, v2 LValue) (LValue, error) {
	f1, err := convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value - f2.Value,
	}, nil
}

func opMul(v1, v2 LValue) (LValue, error) {
	f1, err := convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value * f2.Value,
	}, nil
}

func opDiv(v1, v2 LValue) (LValue, error) {
	f1, err := convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: f1.Value / f2.Value,
	}, nil
}

func opPow(v1, v2 LValue) (LValue, error) {
	f1, err := convToLNumber(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convToLNumber(v2)
	if err != nil {
		return nil, err
	}
	return &LNumber{
		Value: math.Pow(f1.Value, f2.Value),
	}, nil
}
