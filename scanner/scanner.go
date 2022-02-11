// Package scanner
package scanner

import (
	"bufio"
	"os"

	"github.com/jackie8tao/golua/token"
)

const bufSize = 512

type Scanner struct {
	rd   *bufio.Reader
	buf  []rune
	size int
	ch   rune
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
	err = s.next()
	if err != nil {
		return
	}

	for {
		switch s.ch {
		case '+', '-', '*', '%', '^', '#', '&', '|', '(', ')', '[', ']', '{', '}', ',', ';':
			t, err = token.LookupOperator(string(s.ch))
			return
		case '/': // maybe '/' or '//'
			if err = s.next(); err != nil {
				return
			}

			if s.ch == '/' {
				t = token.FLOOR
			} else {
				t = token.DIV
			}
			return
		case '<': // maybe '<=' or '<<' or '<'
			if err = s.next(); err != nil {
				return
			}
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
			if err = s.next(); err != nil {
				return
			}
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
			if err = s.next(); err != nil {
				return
			}
			if s.ch == '=' {
				t = token.EQ
			} else {
				t = token.ASSIGN
			}
			return
		case '~': // maybe '~' or '~='
		case ':': // maybe ':' or '::'
		case '.': // maybe '.' or '..' or '...'
		default:
		}
	}
}

func (s *Scanner) resetBuf() {
	s.buf = make([]rune, bufSize)
	s.size = 0
}

func (s *Scanner) next() error {
	ch, _, err := s.rd.ReadRune()
	if err != nil {
		return err
	}

	s.ch = ch
	return nil
}
