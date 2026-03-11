package ast

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

type Token struct {
	Type TokenType
	Pos  Position
	Name string
	Str  string
}

func (t *Token) String() string {
	return fmt.Sprintf("<type:%v, str:\"%v\">", t.Name, t.Str)
}
