package ast

type Node interface{}

type Stmt interface {
	Node
	stmtNode() //
}

type BaseStmt struct{}

func (s *BaseStmt) stmtNode() {}

type Expr interface {
	Node
	exprNode()
}

type BaseExpr struct{}

func (e *BaseExpr) exprNode() {}
