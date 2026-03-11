package ast

type NilExpr struct{ BaseExpr }

type BoolExpr struct {
	BaseExpr
	Value bool
}

type NumExpr struct {
	BaseExpr
	Value float64
}

type StrExpr struct {
	BaseExpr
	Value string
}

type IdentExpr struct {
	BaseExpr
	Name string
}

// expr -> function
// function -> 'function' funcbody
type FuncExpr struct {
	BaseExpr
	Params     []string
	HasVariant bool // variant params
	Block      *Block
}

// ('[' exp ']' | '.' Name)
// ( args | ':' Name args )
type IndexExpr struct {
	BaseExpr
	Table Expr
	Index Expr
}

type FuncCallExpr struct {
	BaseExpr
	Func Expr
	Args []Expr
}

type MethodCallExpr struct {
	BaseExpr
	Receiver Expr
	Method   string
	Args     []Expr
}

type TableField struct {
	Key   Expr
	Value Expr
}

type TableExpr struct {
	BaseExpr
	Fields []TableField
}

// exp -> exp op exp
type BinOpExpr struct {
	BaseExpr
	LHS Expr
	Op  TokenType
	RHS Expr
}

// exp -> op exp
type UnaryOpExpr struct {
	BaseExpr
	Op  TokenType
	RHS Expr
}
