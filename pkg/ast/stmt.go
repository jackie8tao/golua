package ast

// chunk -> { stat [”;] }
type Chunk struct {
	Stmts []Stmt
}

// block -> chunk
type Block struct {
	Chunk *Chunk
}

// stmt -> varlist1 '=' explist1
type AssignStmt struct {
	BaseStmt
	LHS []Expr
	RHS []Expr
}

// 'local' namelist [init]
type LocalAssignStmt struct {
	BaseStmt
	Names []string
	Exprs []Expr
}

// 'do' block 'end'
type DoStmt struct {
	BaseStmt
	Block *Block
}

// 'while' exp 'do' block 'end'
type WhileStmt struct {
	BaseStmt
	Cond  Expr
	Block *Block
}

// 'repeat' block 'until' exp
type RepeatStmt struct {
	BaseStmt
	Cond  Expr
	Block *Block
}

// 'if' exp 'then' block { elseif exp 'then' block } ['else' block]
type IfStmt struct {
	BaseStmt
	Cond   Expr
	Then   *Block
	ElseIf []*ElseIfSeg
	Else   *Block
}

type ElseIfSeg struct {
	BaseStmt
	Cond Expr
	Then *Block
}

// 'return' [explist]
type ReturnStmt struct {
	BaseStmt
	Exprs []Expr
}

// 'break'
type BreakStmt struct {
	BaseStmt
}

// 'for' Name '=' exp ',' exp [',' exp] 'do' block 'end'
type NumericForStmt struct {
	BaseStmt
	VarName  string
	VarExpr  Expr
	CondExpr Expr
	StepExpr Expr
	Block    *Block
}

// 'for' Name {',' Name} 'in' explist 'do' block 'end'
type GenericForStmt struct {
	BaseStmt
	Names []string
	Exprs []Expr
	Block *Block
}

// 'function' funcname funcbody
type FuncDefStmt struct {
	BaseStmt
	IsLocal    bool // local function definition
	Names      []string
	SuffixName string
	Params     []string
	HasVariant bool // variant params
	Block      *Block
}

// functioncall
type FuncCallStmt struct {
	BaseStmt
	Expr Expr
}
