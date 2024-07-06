package main

type Expr interface {
	Accept(ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinary(expr Binary) interface{}
	VisitGrouping(expr Grouping) interface{}
	VisitLiteral(expr Literal) interface{}
	VisitUnary(expr Unary) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinary(expr)
}

type Grouping struct {
	Expression Expr
}

func (expr Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGrouping(expr)
}

type Literal struct {
	Value interface{}
}

func (expr Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteral(expr)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (expr Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnary(expr)
}
