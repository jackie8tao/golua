package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/jackie8tao/golua/pkg/ast"
	"github.com/jackie8tao/golua/pkg/lexer"
)

// newLexerCmd creates the lexer debugger subcommand.
func newLexerCmd(dir *string) *cobra.Command {
	return newLuaDirCommand("lexer", "Scan .lua files and print tokens", dir, debugLexer)
}

func debugLexer(out io.Writer, path string, r io.Reader) error {
	lx := lexer.NewLexer(r, path)
	for {
		token, err := lx.Scan()
		if err != nil {
			return err
		}
		if token.Type == ast.TokenEOF {
			break
		}
		fmt.Fprintln(out, path, ":", token.String())
	}
	return nil
}
