package main

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/jackie8tao/golua/pkg/lexer"
	"github.com/jackie8tao/golua/pkg/libs"
	"github.com/jackie8tao/golua/pkg/parser"
	"github.com/jackie8tao/golua/pkg/vm"
)

// newVMCmd creates the VM debugger subcommand.
func newVMCmd(dir *string) *cobra.Command {
	return newLuaDirCommand("vm", "Compile and execute .lua files on the VM", dir, debugVM)
}

func debugVM(_ io.Writer, path string, r io.Reader) error {
	lx := lexer.NewLexer(r, path)
	ps := parser.NewParser(lx)
	block, err := ps.Parse()
	if err != nil {
		return err
	}

	c := vm.NewCompiler()
	funcProto, err := c.Compile(block)
	if err != nil {
		return err
	}

	v := vm.NewVm(funcProto)
	libs.SetupBasicLib(v)
	return v.Execute()
}
