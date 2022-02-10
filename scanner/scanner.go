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
		case '+':
			t = token.ADD
			return
		case '-':
			t = token.SUB
			return
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
