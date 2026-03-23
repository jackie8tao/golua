package libs

import (
	"fmt"

	"github.com/jackie8tao/golua/pkg/vm"
)

var basicHandlers = map[string]vm.LGFunction{
	"print": basicPrint,
}

func SetupBasicLib(v *vm.Vm) {
	v.RegisterModule("", basicHandlers)
}

func basicPrint(v *vm.Vm, argc int) int {
	values := make([]vm.LValue, 0)
	for i := 0; i < argc; i++ {
		values = append(values, v.StPop())
	}
	for _, val := range values {
		fmt.Println(val.String())
	}
	return 0
}
