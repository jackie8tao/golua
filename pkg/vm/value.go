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
	Value float64
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

type LTable struct {
	Metadata LValue
	array    []LValue
	dict     map[LValue]LValue
}

func (l *LTable) String() string {
	return fmt.Sprintf("table: %p", l)
}

func (l *LTable) Type() LValueType {
	return LTTable
}

type LGlobal struct {
}

type LState struct {
	G   *LGlobal
	Env *LTable
}

func (l *LState) String() string {
	return fmt.Sprintf("thread: %p", l)
}

func (l *LState) Type() LValueType {
	return LTThread
}

type LUserData struct {
	Metadata *LValue
	Env      *LValue
}

func (l *LUserData) String() string {
	return fmt.Sprintf("userdata: %p", l)
}

func (l *LUserData) Type() LValueType {
	return LTUserData
}

type LFunction struct {
}

func (l *LFunction) String() string {
	return fmt.Sprintf("function: %p", l)
}

func (l *LFunction) Type() LValueType {
	return LTFunction
}
