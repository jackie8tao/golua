// Package scanner
package scanner

import "github.com/jackie8tao/golua/token"

type Scanner struct {
	file *token.File
}

func (s *Scanner) Scan() (pos token.Pos, token token.Token) {
	return
}
