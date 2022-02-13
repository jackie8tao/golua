// Package scanner
package scanner

import (
	"bufio"
	"os"

	"github.com/jackie8tao/golua/token"
)

const bufSize = 512

type Scanner struct {
	err  error
	rd   *bufio.Reader
	buf  []rune
	size int
	ch   rune
}

func isAlpha(ch rune) bool {
	return false
}

func New(file string) (*Scanner, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &Scanner{
		rd:   bufio.NewReader(fp),
		buf:  make([]rune, bufSize),
		size: 0,
		ch:   0,
	}, nil
}

func (s *Scanner) GetIdentifier() string {
	str := ""
	for _, v := range s.buf[:s.size] {
		str += string(v)
	}

	return str
}

func (s *Scanner) Scan() (t token.Token, err error) {
	s.next()
	if s.hasErr() {
		return 0, s.err
	}

	for {
		if s.hasErr() {
			return 0, s.err
		}

		switch s.ch {
		case '+', '-', '*', '%', '^', '#', '&', '|', '(', ')', '[', ']', '{', '}', ',', ';':
			t, err = token.LookupOperator(string(s.ch))
			return
		case '/': // maybe '/' or '//'
			s.next()
			if s.ch == '/' {
				t = token.FLOOR
				return
			}
			t = token.DIV
			return
		case '<': // maybe '<=' or '<<' or '<'
			s.next()
			switch s.ch {
			case '<':
				t = token.BT_SHL
			case '=':
				t = token.LEQ
			default:
				t = token.LT
			}
			return
		case '>': // maybe '>=' or '>>' or '>'
			s.next()
			switch s.ch {
			case '=':
				t = token.GEQ
			case '>':
				t = token.BT_SHR
			default:
				t = token.GT
			}
			return
		case '=': // maybe '=' or '=='
			s.next()
			if s.ch == '=' {
				t = token.EQ
				return
			}
			t = token.ASSIGN
			return
		case '~': // maybe '~' or '~='
			s.next()
			if s.ch == '=' {
				t = token.NEQ
				return
			}
			t = token.BT_XOR
			return
		case ':': // maybe ':' or '::'
			s.next()
			if s.ch == ':' {
				t = token.DBCOLON
				return
			}
			t = token.COLON
			return
		case '.': // maybe '.' or '..' or '...'
			s.next()
			if s.ch != '.' {
				t = token.DOT
				return
			}
			s.next()
			if s.ch != '.' {
				t = token.CONCAT
				return
			}
			t = token.DOTS
			return
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

		default:
			if isAlpha(s.ch) {
				return
			}
		}
	}
}

func (s *Scanner) hasErr() bool {
	return s.err != nil
}

func (s *Scanner) resetBuf() {
	s.buf = make([]rune, bufSize)
	s.size = 0
}

func (s *Scanner) next() {
	ch, _, err := s.rd.ReadRune()
	if err != nil {
		s.err = err
		return
	}
	s.ch = ch
	return
}
