package ast

type Node interface{}

type Stmt interface {
	Node
	stmtNode() // 占位方法，用于让 Go 编译器在编译期做严格类型约束
}

type BaseStmt struct{}

func (s *BaseStmt) stmtNode() {}

type Expr interface {
	Node
	exprNode() // 占位方法
}

type BaseExpr struct{}

func (e *BaseExpr) exprNode() {}
