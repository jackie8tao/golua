package lexer

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenString
	TokenAnd       // and
	TokenBreak     // break
	TokenDo        // do
	TokenElse      // else
	TokenElseif    // elseif
	TokenEnd       // end
	TokenFalse     // false
	TokenFor       // for
	TokenFunction  // function
	TokenIf        // if
	TokenIn        // in
	TokenLocal     // local
	TokenNil       // nil
	TokenNot       // not
	TokenOr        // or
	TokenRepeat    // repeat
	TokenReturn    // return
	TokenThen      // then
	TokenTrue      // true
	TokenUntil     // until
	TokenWhile     // while
	TokenAdd       // +
	TokenMinus     // -
	TokenMultiply  // *
	TokenDivide    // /
	TokenXor       // ^
	TokenAssign    // =
	TokenNeq       // ~=
	TokenLeq       // <=
	TokenGeq       // >=
	TokenLt        // <
	TokenGt        // >
	TokenEq        // ==
	TokenLparen    // (
	TokenRparen    // )
	TokenLbracket  // [
	TokenRbracket  // ]
	TokenLbrace    // {
	TokenRbrace    // }
	TokenSemicolon // ;
	TokenComma     // ,
	TokenColon     // :
	TokenDot       // .
	TokenDotDot    // ..
	TokenDots      // ...
)

type Position struct {
	Source string
	Line   int
	Column int
}

var reservedWords = map[string]TokenType{
	"and":      TokenAnd,
	"break":    TokenBreak,
	"do":       TokenDo,
	"else":     TokenElse,
	"elseif":   TokenElseif,
	"end":      TokenEnd,
	"false":    TokenFalse,
	"for":      TokenFor,
	"function": TokenFunction,
	"if":       TokenIf,
	"in":       TokenIn,
	"local":    TokenLocal,
	"nil":      TokenNil,
	"not":      TokenNot,
	"or":       TokenOr,
	"repeat":   TokenRepeat,
	"return":   TokenReturn,
	"then":     TokenThen,
	"true":     TokenTrue,
	"until":    TokenUntil,
	"while":    TokenWhile,
}

var operators = map[string]TokenType{
	"+":   TokenAdd,
	"-":   TokenMinus,
	"*":   TokenMultiply,
	"/":   TokenDivide,
	"^":   TokenXor,
	"=":   TokenAssign,
	"~=":  TokenNeq,
	"<=":  TokenLeq,
	">=":  TokenGeq,
	"<":   TokenLt,
	">":   TokenGt,
	"==":  TokenEq,
	"(":   TokenLparen,
	")":   TokenRparen,
	"[":   TokenLbracket,
	"]":   TokenRbracket,
	"{":   TokenLbrace,
	"}":   TokenRbrace,
	";":   TokenSemicolon,
	",":   TokenComma,
	":":   TokenColon,
	".":   TokenDot,
	"..":  TokenDotDot,
	"...": TokenDots,
}

var tokenNames = map[TokenType]string{
	TokenEOF:        "TK_EOF",
	TokenIdentifier: "TK_IDENTIFIER",
	TokenNumber:     "TK_NUMBER",
	TokenString:     "TK_STRING",
	TokenAnd:        "TK_AND",
	TokenBreak:      "TK_BREAK",
	TokenDo:         "TK_DO",
	TokenElse:       "TK_ELSE",
	TokenElseif:     "TK_ELSEIF",
	TokenEnd:        "TK_END",
	TokenFalse:      "TK_FALSE",
	TokenFor:        "TK_FOR",
	TokenFunction:   "TK_FUNCTION",
	TokenIf:         "TK_IF",
	TokenIn:         "TK_IN",
	TokenLocal:      "TK_LOCAL",
	TokenNil:        "TK_NIL",
	TokenNot:        "TK_NOT",
	TokenOr:         "TK_OR",
	TokenRepeat:     "TK_REPEAT",
	TokenReturn:     "TK_RETURN",
	TokenThen:       "TK_THEN",
	TokenTrue:       "TK_TRUE",
	TokenUntil:      "TK_UNTIL",
	TokenWhile:      "TK_WHILE",
	TokenAdd:        "TK_ADD",
	TokenMinus:      "TK_MINUS",
	TokenMultiply:   "TK_MULTIPLY",
	TokenDivide:     "TK_DIVIDE",
	TokenXor:        "TK_XOR",
	TokenAssign:     "TK_ASSIGN",
	TokenNeq:        "TK_NEQ",
	TokenLeq:        "TK_LEQ",
	TokenGeq:        "TK_GEQ",
	TokenLt:         "TK_LT",
	TokenGt:         "TK_GT",
	TokenEq:         "TK_EQ",
	TokenLparen:     "TK_LPAREN",
	TokenRparen:     "TK_RPAREN",
	TokenLbracket:   "TK_LBRACKET",
	TokenRbracket:   "TK_RBRACKET",
	TokenLbrace:     "TK_LBRACE",
	TokenRbrace:     "TK_RBRACE",
	TokenSemicolon:  "TK_SEMICOLON",
	TokenComma:      "TK_COMMA",
	TokenColon:      "TK_COLON",
	TokenDot:        "TK_DOT",
	TokenDotDot:     "TK_DOTDOT",
	TokenDots:       "TK_DOTS",
}

type Token struct {
	Type TokenType
	Pos  Position
	Name string
	Str  string
}

func (t *Token) String() string {
	return fmt.Sprintf("<type:%v, str:%v>", t.Name, t.Str)
}
