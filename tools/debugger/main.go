package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jackie8tao/golua/pkg/ast"
	"github.com/jackie8tao/golua/pkg/lexer"
	"github.com/jackie8tao/golua/pkg/libs"
	"github.com/jackie8tao/golua/pkg/parser"
	"github.com/jackie8tao/golua/pkg/vm"
)

type luaFileHandler func(out io.Writer, path string, r io.Reader) error

type subcommandSpec struct {
	use     string
	short   string
	handler luaFileHandler
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

// newRootCmd builds the root debugger command and registers all subcommands.
func newRootCmd() *cobra.Command {
	var dir string

	rootCmd := &cobra.Command{
		Use:   "debugger",
		Short: "Unified CLI debugger for GoLua lexer, parser, and VM",
	}
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "directory to scan for .lua files")

	specs := []subcommandSpec{
		{
			use:     "lexer",
			short:   "Scan .lua files and print tokens",
			handler: debugLexer,
		},
		{
			use:     "parser",
			short:   "Parse .lua files and print AST blocks",
			handler: debugParser,
		},
		{
			use:     "vm",
			short:   "Compile and execute .lua files on the VM",
			handler: debugVM,
		},
	}

	for _, spec := range specs {
		s := spec
		cmd := &cobra.Command{
			Use:   s.use,
			Short: s.short,
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runOnLuaDir(cmd.OutOrStdout(), dir, s.handler)
			},
		}
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// runOnLuaDir iterates over .lua files in dir and runs handler on each file.
func runOnLuaDir(out io.Writer, dir string, handler luaFileHandler) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".lua" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := runOnLuaFile(out, path, handler); err != nil {
			// Report per-file failures and continue inspecting remaining files.
			fmt.Fprintln(out, path, ":", err)
		}
	}
	return nil
}

func runOnLuaFile(out io.Writer, path string, handler luaFileHandler) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	return handler(out, path, fp)
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
