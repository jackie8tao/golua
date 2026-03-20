package ast

// Chunk chunk -> { stat [”;] }
type Chunk struct {
	Stmts []Stmt
}

// Block block -> chunk
type Block struct {
	Chunk *Chunk
}

// AssignStmt stmt -> varlist1 '=' explist1
type AssignStmt struct {
	BaseStmt
	LHS []Expr
	RHS []Expr
}

// LocalAssignStmt 'local' namelist [init]
type LocalAssignStmt struct {
	BaseStmt
	Names []string
	Exprs []Expr
}

// DoStmt 'do' block 'end'
type DoStmt struct {
	BaseStmt
	Block *Block
}

// WhileStmt 'while' exp 'do' block 'end'
type WhileStmt struct {
	BaseStmt
	Cond  Expr
	Block *Block
}

// RepeatStmt 'repeat' block 'until' exp
type RepeatStmt struct {
	BaseStmt
	Cond  Expr
	Block *Block
}

// IfStmt 'if' exp 'then' block { elseif exp 'then' block } ['else' block]
type IfStmt struct {
	BaseStmt
	Cond   Expr
	Then   *Block
	ElseIf []*ElseIfSeg
	Else   *Block
}

type ElseIfSeg struct {
	Cond Expr
	Then *Block
}

// ReturnStmt 'return' [explist]
type ReturnStmt struct {
	BaseStmt
	Exprs []Expr
}

// BreakStmt 'break'
type BreakStmt struct {
	BaseStmt
}

// NumericForStmt 'for' Name '=' exp ',' exp [',' exp] 'do' block 'end'
type NumericForStmt struct {
	BaseStmt
	VarName  string
	VarExpr  Expr
	CondExpr Expr
	StepExpr Expr
	Block    *Block
}

// GenericForStmt 'for' Name {',' Name} 'in' explist 'do' block 'end'
type GenericForStmt struct {
	BaseStmt
	Names []string
	Exprs []Expr
	Block *Block
}

// FuncDefStmt 'function' funcname funcbody
type FuncDefStmt struct {
	BaseStmt
	IsLocal    bool // local function definition
	Names      []string
	SuffixName string
	Params     []string
	HasVariant bool // variant params
	Block      *Block
}

// FuncCallStmt functioncall
type FuncCallStmt struct {
	BaseStmt
	Expr Expr
}
