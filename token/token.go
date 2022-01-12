// Package token defines consts representing the lexical tokens of the Lua
// programming language and basic operations on tokens (printing, predicates).
package token

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	IDENT
	INT
	FLOAT
	IMAG
	CHAR
	STRING
	literal_end

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

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	operator_end
)
