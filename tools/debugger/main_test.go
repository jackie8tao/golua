package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLuaFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file %s: %v", path, err)
	}
	return path
}

func runDebuggerCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestDebuggerSubcommands(t *testing.T) {
	dir := t.TempDir()
	luaPath := writeTempLuaFile(t, dir, "sample.lua", "local x = 1 + 2\n")
	// Non-.lua files and subdirectories must be ignored.
	writeTempLuaFile(t, dir, "notes.txt", "not lua")
	if err := os.Mkdir(filepath.Join(dir, "nested.lua"), 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	t.Run("lexer", func(t *testing.T) {
		out, err := runDebuggerCmd(t, "lexer", "--dir", dir)
		if err != nil {
			t.Fatalf("lexer command failed: %v", err)
		}
		if !strings.Contains(out, luaPath) || !strings.Contains(out, "TK_LOCAL") {
			t.Fatalf("unexpected lexer output: %q", out)
		}
		if strings.Contains(out, "notes.txt") || strings.Contains(out, "nested.lua") {
			t.Fatalf("non-file or non-.lua entry should be skipped, got: %q", out)
		}
	})

	t.Run("parser", func(t *testing.T) {
		out, err := runDebuggerCmd(t, "parser", "-d", dir)
		if err != nil {
			t.Fatalf("parser command failed: %v", err)
		}
		if !strings.Contains(out, luaPath) {
			t.Fatalf("unexpected parser output: %q", out)
		}
	})

	t.Run("vm valid script", func(t *testing.T) {
		out, err := runDebuggerCmd(t, "vm", "--dir", dir)
		if err != nil {
			t.Fatalf("vm command failed: %v", err)
		}
		if strings.TrimSpace(out) != "" {
			t.Fatalf("expected no error output for valid script, got: %q", out)
		}
	})

	t.Run("vm reports per-file error without aborting", func(t *testing.T) {
		errDir := t.TempDir()
		badPath := writeTempLuaFile(t, errDir, "bad.lua", "return 1\n")
		out, err := runDebuggerCmd(t, "vm", "--dir", errDir)
		if err != nil {
			t.Fatalf("vm command should not fail on per-file error, got: %v", err)
		}
		if !strings.Contains(out, badPath) || !strings.Contains(out, "invalid ast node") {
			t.Fatalf("expected per-file error in output, got: %q", out)
		}
	})

	t.Run("missing directory returns error", func(t *testing.T) {
		missing := filepath.Join(dir, "does-not-exist")
		if _, err := runDebuggerCmd(t, "lexer", "--dir", missing); err == nil {
			t.Fatal("expected error for missing directory, got nil")
		}
	})
}
