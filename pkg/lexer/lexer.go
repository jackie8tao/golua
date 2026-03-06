package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

const runeEOF = -1

type Error struct {
	Pos     Position
	Message string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == runeEOF {
		return fmt.Sprintf("%v at EOF:   %s\n", pos.Source, e.Message)
	}
	return fmt.Sprintf(
		"%v line:%d(column:%d) : %s\n",
		pos.Source, pos.Line, pos.Column, e.Message,
	)
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

func writeRune(buf *bytes.Buffer, ch rune) {
	buf.WriteRune(ch)
}

func (l *Lexer) newError(msg string) *Error {
	return &Error{
		Pos:     l.Pos,
		Message: msg,
	}
}

func (l *Lexer) readNext() rune {
	ch, _, err := l.reader.ReadRune()
	if err == io.EOF {
		return runeEOF
	}
	return ch
}

func (l *Lexer) newline(ch rune) {
	if ch < 0 {
		return
	}
	l.Pos.Line++
	l.Pos.Column = 0
	next := l.peek()
	if (ch == '\n' && next == '\r') || (ch == '\r' && next == '\n') {
		_, _, _ = l.reader.ReadRune()
	}
}

func (l *Lexer) next() rune {
	ch := l.readNext()
	switch ch {
	case '\n', '\r':
		l.newline(ch)
		ch = '\n'
	case runeEOF:
		l.Pos.Line = runeEOF
		l.Pos.Column = 0
	default:
		l.Pos.Column++
	}
	return ch
}

func (l *Lexer) peek() rune {
	ch := l.readNext()
	if ch != runeEOF {
		_ = l.reader.UnreadRune()
	}
	return ch
}

func (l *Lexer) scanNewline() {
	oldCh := l.next() // \n or \r

	// windows \r\n or \n\r
	newCh := l.peek()
	if newCh == '\n' || newCh == '\r' {
		if oldCh != newCh {
			l.next()
		}
	}
}

func (l *Lexer) scanDecimal(buf *bytes.Buffer) error {
	writeRune(buf, l.next())
	hasPoint := false
	hasExponent := false
redo:
	ch := l.peek()
	switch {
	case isDecimal(ch): // decimal
		writeRune(buf, l.next())
		goto redo
	default:
		switch ch {
		case 'e', 'E': // exponent
			if hasExponent {
				return l.newError("unexpected exponent")
			}
			hasExponent = true
			writeRune(buf, l.next())

			ch = l.peek()
			// optional sign
			if ch == '+' || ch == '-' {
				writeRune(buf, l.next())
				ch = l.peek()
			}
			// exponent first digits
			if !isDecimal(ch) {
				return l.newError("invalid exponent")
			}
			writeRune(buf, l.next())
			goto redo
		case '.': // point
			if hasPoint {
				return l.newError("unexpected point")
			}
			hasPoint = true
			writeRune(buf, l.next())
			goto redo
		default:
			return nil
		}
	}
}

func (l *Lexer) scanIdentifier(buf *bytes.Buffer) {
	writeRune(buf, l.next())
	for {
		ch := l.peek()
		if !isIdent(ch, 1) {
			break
		}
		writeRune(buf, l.next())
	}
}

func (l *Lexer) scanEscape(buf *bytes.Buffer) error {
	ch := l.next()
	switch ch {
	case 'a':
		buf.WriteRune('\a')
	case 'b':
		buf.WriteRune('\b')
	case 'f':
		buf.WriteRune('\f')
	case 'n':
		buf.WriteRune('\n')
	case 'r':
		buf.WriteRune('\r')
	case 't':
		buf.WriteRune('\t')
	case 'v':
		buf.WriteRune('\v')
	case '\\':
		buf.WriteRune('\\')
	case '"':
		buf.WriteRune('"')
	case '\'':
		buf.WriteRune('\'')
	case '\n':
		buf.WriteRune('\n')
	case '\r':
		buf.WriteRune('\n')
	default:
		if ch >= '0' && ch <= '9' {
			bs := []rune{ch}
			for i := 0; i < 2 && isDecimal(l.peek()); i++ {
				bs = append(bs, l.next())
			}
			val, err := strconv.ParseInt(string(bs), 10, 32)
			if err != nil {
				return l.newError("invalid decimal escape sequence")
			}
			buf.WriteRune(rune(val))
		} else {
			buf.WriteRune(ch)
		}
	}

	return nil
}

func (l *Lexer) scanString(buf *bytes.Buffer) error {
	quote := l.next()
	for {
		ch := l.next()
		if ch == runeEOF || ch == '\n' || ch == '\r' {
			return l.newError("unterminated string")
		}
		if ch == quote {
			break
		}
		if ch == '\\' {
			if err := l.scanEscape(buf); err != nil {
				return err
			}
		} else {
			writeRune(buf, ch)
		}
	}
	return nil
}

func (l *Lexer) scanMultilineString(buf *bytes.Buffer) error {
	ch := l.next()
	if ch == '\n' || ch == '\r' {
		ch = l.next()
	}
	for {
		ch = l.next()
		if ch == runeEOF {
			return l.newError("unterminated multiline string")
		}
		if ch == ']' && l.peek() == ']' {
			_ = l.next()
			break
		}
		writeRune(buf, ch)
	}
	return nil
}

func (l *Lexer) skipComments() error {
	// multiple line comment
	ch := l.next()
	if ch == '[' && l.next() == '[' {
		if err := l.scanMultilineString(&bytes.Buffer{}); err != nil {
			return l.newError("invalid multiline comment")
		}
	}

	// single line comment
	for {
		ch = l.next()
		if ch == '\n' || ch == '\r' {
			break
		}
	}

	return nil
}

func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

func isIdent(ch rune, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || (isDecimal(ch) && pos > 0)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNewline(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func (l *Lexer) Scan() (token Token, err error) {
redo:
	ch := l.peek()
	buf := &bytes.Buffer{}

	switch {
	case isWhitespace(ch): // whitespace
		l.next()
		goto redo
	case isNewline(ch): // newline
		l.scanNewline()
		goto redo
	case isDecimal(ch):
		err = l.scanDecimal(buf)
		if err != nil {
			return
		}
		token.Type = TokenNumber
		token.Str = buf.String()
		goto finally
	case isIdent(ch, 0):
		l.scanIdentifier(buf)
		token.Type = TokenIdentifier
		if tokenType, ok := reservedWords[buf.String()]; ok {
			token.Type = tokenType
		}
		token.Str = buf.String()
		goto finally
	default:
		switch ch {
		case runeEOF:
			token.Type = TokenEOF
			goto finally
		case '+', '*', '/', '%', '^', '#', '(', ')', '{', '}', ']', ';', ',', ':':
			writeRune(buf, l.next())
			token.Str = buf.String()
			token.Type = operators[token.Str]
			goto finally
		case '[':
			writeRune(buf, l.next())
			if ch = l.peek(); ch == '[' {
				token.Type = TokenString
				buf.Reset()
				if err = l.scanMultilineString(buf); err != nil {
					return
				}
			} else {
				token.Type = TokenLbracket
			}
			token.Str = buf.String()
			goto finally
		case '~':
			writeRune(buf, l.next())
			ch = l.peek()
			if ch == '=' {
				token.Type = TokenNeq
				token.Str = buf.String()
				goto finally
			}
			err = l.newError("invalid token")
			return
		case '-':
			writeRune(buf, l.next())
			token.Type = TokenMinus
			if ch = l.peek(); ch == '-' { // comment
				_ = l.next()
				if err = l.skipComments(); err != nil {
					return
				}
				goto redo
			}
			token.Str = buf.String()
			goto finally
		case '=':
			writeRune(buf, l.next())
			token.Type = TokenAssign
			if ch = l.peek(); ch == '=' {
				writeRune(buf, l.next())
				token.Type = TokenEq
			}
			token.Str = buf.String()
			goto finally
		case '<':
			writeRune(buf, l.next())
			token.Type = TokenLt
			if ch = l.peek(); ch == '=' {
				writeRune(buf, l.next())
				token.Type = TokenLeq
			}
			token.Str = buf.String()
			goto finally
		case '>':
			writeRune(buf, l.next())
			token.Type = TokenGt
			if ch = l.peek(); ch == '=' {
				writeRune(buf, l.next())
				token.Type = TokenGeq
			}
			token.Str = buf.String()
			goto finally
		case '.':
			writeRune(buf, l.next())
			token.Type = TokenDot
			if ch = l.peek(); ch == '.' {
				writeRune(buf, l.next())
				token.Type = TokenDotDot
				if ch = l.peek(); ch == '.' {
					writeRune(buf, l.next())
					token.Type = TokenDots
				}
			}
			token.Str = buf.String()
			goto finally
		case '"', '\'':
			err = l.scanString(buf)
			if err != nil {
				return
			}
			token.Type = TokenString
			token.Str = buf.String()
			goto finally
		default:
			err = l.newError("invalid token")
			return
		}
	}

finally:
	token.Pos = l.Pos
	token.Name = tokenNames[token.Type]
	return
}
