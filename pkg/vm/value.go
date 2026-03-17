package vm

import "fmt"

type LValueType int

const (
	LTNil LValueType = iota
	LTBool
	LTNumber
	LTString
	LTFunction
	LTUserData
	LTThread
	LTTable
)

var lValueNames = [...]string{
	"nil",
	"boolean",
	"number",
	"string",
	"function",
	"userdata",
	"thread",
	"table",
}

func (vt LValueType) String() string {
	return lValueNames[int(vt)]
}

type LValue interface {
	String() string
	Type() LValueType
}

type LNil struct{}

func (l *LNil) String() string {
	return "nil"
}

func (l *LNil) Type() LValueType {
	return LTNil
}

type LBool struct {
	Value bool
}

func (l *LBool) String() string {
	if l.Value {
		return "true"
	}
	return "false"
}

func (l *LBool) Type() LValueType {
	return LTBool
}

type LNumber struct {
	Value float32
}

func (l *LNumber) String() string {
	return fmt.Sprintf("%.2f", l.Value)
}

func (l *LNumber) Type() LValueType {
	return LTNumber
}

type LString struct {
	Value string
}

func (l *LString) String() string {
	return l.Value
}

func (l *LString) Type() LValueType {
	return LTString

}
