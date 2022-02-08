// Package token defines consts representing the lexical tokens of the Lua
// programming language and basic operations on tokens (printing, predicates).
package token

// Token is the set of lexical tokens of lua programming language.
type Token int

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
	AND1 // &
	OR1  // |
	XOR  // ~
	SHR  // >>
	SHL  // <<

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
	LEN // #

	// delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACK    // [
	RBRACK    // ]
	LBRACE    // {
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

var tokens = []string{
	ADD:      "+",
	AND:      "and",
	BREAK:    "break",
	DO:       "do",
	ELSE:     "else",
	ELSEIF:   "elseif",
	END:      "end",
	FALSE:    "false",
	FOR:      "for",
	FUNCTION: "function",
	GOTO:     "goto",
	IF:       "if",
	IN:       "in",
	LOCAL:    "local",
	NIL:      "nil",
	NOT:      "not",
	OR:       "or",
	REPEAT:   "repeat",
	RETURN:   "return",
	THEN:     "then",
	TRUE:     "true",
	UNTIL:    "until",
	WHILE:    "while",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}
