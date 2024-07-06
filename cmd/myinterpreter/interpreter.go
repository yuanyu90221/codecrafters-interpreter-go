package main

import (
	"fmt"
	"strings"
)

type Interpreter struct{}

func (i *Interpreter) interpret(expr Expr) {
	value := i.Evaluate(expr)
	fmt.Println(stringify(value))
}

func (i *Interpreter) Evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinary(expr Binary) interface{} {
	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Type {
	case MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case PLUS:
		if lVal, ok := left.(float64); ok {
			if rVal, ok := right.(float64); ok {
				return lVal + rVal
			}
		}
		if lVal, ok := left.(string); ok {
			if rVal, ok := right.(string); ok {
				return lVal + rVal
			}
		}
		panic(fmt.Sprintf("%s operand must be a two numbers or two strings.", expr.Operator.Lexeme))
	case GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}
	return nil
}
func (i *Interpreter) VisitLiteral(expr Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnary(expr Unary) interface{} {
	right := i.Evaluate(expr.Right)
	switch expr.Operator.Type {
	case MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case BANG:
		return !isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitGrouping(expr Grouping) interface{} {
	return i.Evaluate(expr.Expression)
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if object == false {
		return false
	}
	return true
}
func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if _, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
		return text
	}
	return fmt.Sprintf("%v", object)
}
func checkNumberOperand(operator Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(fmt.Sprintf("%s operand must be a number.", operator.Lexeme))
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(fmt.Sprintf("%s operand must be a number.", operator.Lexeme))
}
