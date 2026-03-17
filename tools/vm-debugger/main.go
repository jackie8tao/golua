package main

import "github.com/jackie8tao/golua/pkg/vm"

func main() {
	v := vm.NewVm()
	v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LNumber{Value: 2}))
	v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LNumber{Value: 1}))
	v.WriteCode(vm.OpPow)
	v.WriteCode(vm.OpPrint)
	v.WriteCode(vm.OpConstant, v.WriteConstant(&vm.LString{Value: "Hello World"}))
	v.WriteCode(vm.OpPrint)
	err := v.Execute()
	if err != nil {
		panic(err)
	}
}
