// Package token defines consts representing the lexical tokens of the Lua
// programming language and basic operations on tokens (printing, predicates).
package token

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT
)
