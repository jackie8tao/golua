package main

import "github.com/jackie8tao/golua/pkg/vm"

func main() {
	v := vm.NewVm()
	//v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LNumber{Value: 2}))
	//v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LNumber{Value: 1}))
	//v.WriteCode(vm.OpPow)
	//v.WriteCode(vm.OpPrint)
	//v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LString{Value: "Hello World"}))
	//v.WriteCode(vm.OpPrint)
	v.WriteGlobals("a", &vm.LNumber{Value: 1})
	v.WriteGlobals("b", &vm.LNumber{Value: 2})
	aIdx := v.WriteConstant(&vm.LString{Value: "a"})
	bIdx := v.WriteConstant(&vm.LString{Value: "b"})
	v.WriteCode(vm.OpGetGlobal, aIdx)
	v.WriteCode(vm.OpGetGlobal, bIdx)
	v.WriteCode(vm.OpAdd)
	v.WriteCode(vm.OpSetGlobal, aIdx)
	v.WriteCode(vm.OpGetGlobal, aIdx)
	v.WriteCode(vm.OpPrint)
	err := v.Execute()
	if err != nil {
		panic(err)
	}
}
