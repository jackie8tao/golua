// Package scanner
package scanner

import (
	"os"

	"github.com/jackie8tao/golua/token"
)

type Scanner struct {
	fp  *os.File
	buf []byte
	ch  rune
}

func New(file string) *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan() (pos token.Pos, token token.Token) {
	switch s.ch {
	case '+':
	}

	return
}
