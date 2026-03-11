package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackie8tao/golua/pkg/lexer"
	"github.com/jackie8tao/golua/pkg/parser"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	flag.Parse()

	entries, err := os.ReadDir(*dir)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if filepath.Ext(v.Name()) != ".lua" {
			continue
		}
		path := filepath.Join(*dir, v.Name())
		if err = parseLua(path); err != nil {
			fmt.Println(path, ":", err)
		}
	}
}

func parseLua(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	lx := lexer.NewLexer(fp, path)
	ps := parser.NewParser(lx)
	block, err := ps.Parse()
	if err != nil {
		return err
	}
	fmt.Println(path, ":", block)
	return nil
}
