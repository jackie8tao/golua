package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/jackie8tao/golua/pkg/lexer"
	"github.com/jackie8tao/golua/pkg/parser"
)

// newParserCmd creates the parser debugger subcommand.
func newParserCmd(dir *string) *cobra.Command {
	return newLuaDirCommand("parser", "Parse .lua files and print AST blocks", dir, debugParser)
}

func debugParser(out io.Writer, path string, r io.Reader) error {
	lx := lexer.NewLexer(r, path)
	ps := parser.NewParser(lx)
	block, err := ps.Parse()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, path, ":", block)
	return nil
}
