package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type luaFileHandler func(out io.Writer, path string, r io.Reader) error

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

	rootCmd.AddCommand(
		newLexerCmd(&dir),
		newParserCmd(&dir),
		newVMCmd(&dir),
	)

	return rootCmd
}

// newLuaDirCommand constructs a subcommand that applies handler to every .lua file in *dir.
func newLuaDirCommand(use, short string, dir *string, handler luaFileHandler) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOnLuaDir(cmd.OutOrStdout(), *dir, handler)
		},
	}
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
