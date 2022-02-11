package token

import "errors"

// error list
var (
	ErrKeyword  = errors.New("invalid keyword")
	ErrOperator = errors.New("invalid operator")
)
