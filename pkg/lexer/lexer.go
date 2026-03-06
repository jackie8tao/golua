package lexer

import (
	"bufio"
	"fmt"
	"io"
)

const EOF = -1

var whitespaces = [...]rune{' ', '\t'}
var newlines = [...]rune{'\n', '\r'}

type Error struct {
	Pos     Position
	Message string
	Token   string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == EOF {
		return fmt.Sprintf("%v at EOF:   %s\n", pos.Source, e.Message)
	}
	return fmt.Sprintf(
		"%v line:%d(column:%d) near '%v': %s\n",
		pos.Source, pos.Line, pos.Column, e.Token, e.Message,
	)
}

func isDecimal(ch int) bool { return '0' <= ch && ch <= '9' }

func isIdent(ch int, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || isDecimal(ch) && pos > 0
}

func isDigit(ch int) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

type Lexer struct {
	Pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader, source string) *Lexer {
	return &Lexer{
		reader: bufio.NewReaderSize(reader, 4096),
		Pos: Position{
			Source: source,
			Line:   1,
			Column: 0,
		},
	}
}

func (l *Lexer) peek() rune {
	peek, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic("")
	}
	l.reader.UnreadRune()
	return peek
}

func (l *Lexer) Next() (Token, error) {
	return Token{}, nil
}
