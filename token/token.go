// Package token defines consts representing the lexical tokens of the Lua
// programming language and basic operations on tokens (printing, predicates).
package token

import (
	"fmt"
)

// Token is the set of lexical tokens of lua programming language.
type Token int

// constant values for token of lua programming language.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifier and basic type literals
	IDENT
	NUMBER
	STRING
	literal_end

	operator_beg
	// Operators and delimiters
	// arithmetic operators
	ADD   // +
	SUB   // -
	MUL   // *
	DIV   // /
	FLOOR // //
	MOD   // %
	EXP   // ^

	// bitwise operators
	BT_AND // &
	BT_OR  // |
	BT_XOR // ~
	BT_SHR // >>
	BT_SHL // <<

	// relational operators
	EQ  // ==
	NEQ // ~=
	LT  // <
	GT  // >
	LEQ // <=
	GEQ // >=

	// concatenate operators
	CONCAT // ..

	// length operators
	LENGTH // #

	// assignment operators
	ASSIGN // =

	// delimiters
	LPAREN    // (
	LBRACK    // [
	LBRACE    // {
	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	COLON     // :
	DBCOLON   // ::
	DOT       // .
	DOTS      // ...
	COMMA     // ,
	SEMICOLON // ;
	operator_end

	keyword_beg
	// Keywords
	AND      // and
	BREAK    // break
	DO       // do
	ELSE     // else
	ELSEIF   // elseif
	END      // end
	FALSE    // false
	FOR      // for
	FUNCTION // function
	GOTO     // goto
	IF       // if
	IN       // in
	LOCAL    // local
	NIL      // nil
	NOT      // not
	OR       // or
	REPEAT   // repeat
	RETURN   // return
	THEN     // then
	TRUE     // true
	UNTIL    // until
	WHILE    // while
	keyword_end
)

// all reversed keywords and symbols of lua5.4
var tokens = []string{
	ADD:       "+",
	SUB:       "-",
	MUL:       "*",
	DIV:       "/",
	FLOOR:     "//",
	MOD:       "%",
	EXP:       "^",
	BT_AND:    "&",
	BT_OR:     "|",
	BT_XOR:    "~",
	BT_SHR:    ">>",
	BT_SHL:    "<<",
	AND:       "and",
	OR:        "or",
	NOT:       "not",
	BREAK:     "break",
	DO:        "do",
	ELSE:      "else",
	ELSEIF:    "elseif",
	END:       "end",
	FALSE:     "false",
	FOR:       "for",
	FUNCTION:  "function",
	GOTO:      "goto",
	IF:        "if",
	IN:        "in",
	LOCAL:     "local",
	NIL:       "nil",
	REPEAT:    "repeat",
	RETURN:    "return",
	THEN:      "then",
	TRUE:      "true",
	UNTIL:     "until",
	WHILE:     "while",
	LPAREN:    "(",
	LBRACK:    "[",
	LBRACE:    "{",
	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	COLON:     ":",
	DBCOLON:   "::",
	DOT:       ".",
	CONCAT:    "..",
	DOTS:      "...",
	COMMA:     ",",
	SEMICOLON: ";",
	ASSIGN:    "=",
}

// keywords and operators of lua5.4
// used to simplify the recognization of tokens.
var (
	keywords  = make(map[string]Token)
	operators = make(map[string]Token)
)

func init() {
	initKeywords()
	initOperators()
}

func initKeywords() {
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func initOperators() {
	for i := operator_beg + 1; i < operator_end; i++ {
		operators[tokens[i]] = i
	}
}

// LookupOperator lookup operator token from operators table.
func LookupOperator(key string) (Token, error) {
	val, ok := operators[key]
	if !ok {
		return 0, ErrOperator
	}
	return val, nil
}

// LookupKeyword lookup keyword token from keywords table.
func LookupKeyword(key string) (Token, error) {
	val, ok := keywords[key]
	if ok {
		return 0, ErrKeyword
	}
	return val, nil
}

// PrintKeywords print lua keywords to terminal.
func PrintKeywords() {
	for _, v := range keywords {
		fmt.Println(v)
	}
}
