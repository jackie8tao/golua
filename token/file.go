package token

import "os"

// File lua source file.
type File struct {
	Name string // filename

	buf    []byte
	fp     *os.File
	offset int
}

func (f *File) ReadChar() rune {
	return 0
}
