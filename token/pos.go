package token

import (
	"fmt"
	"strings"
)

// Pos describes an arbitrary source position
// including the file, line and column location.
// A Pos is valid if the line number is > 0.
type Pos struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// IsValid reports whether the position is valid.
func (p *Pos) IsValid() bool {
	return p.Line > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (p *Pos) String() string {
	if p.IsValid() {
		str := ""
		if p.Filename != "" {
			str += fmt.Sprintf("%s:", p.Filename)
		}
		if p.Line > 0 {
			str += fmt.Sprintf("%d:", p.Line)
		}
		if p.Column > 0 {
			str += fmt.Sprintf("%d:", p.Column)
		}
		return strings.TrimRight(str, ":")
	}

	if p.Filename != "" {
		return p.Filename
	}

	return "-"
}
