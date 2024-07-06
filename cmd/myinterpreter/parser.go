package main

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparsion()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparsion()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparsion() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LEFT_BRACE) {
		operator := p.previous()
		right := p.term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{Operator: operator, Right: right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{Value: false}
	}
	if p.match(TRUE) {
		return Literal{Value: true}
	}
	if p.match(NIL) {
		return Literal{Value: nil}
	}
	if p.match(NUMBER, STRING) {
		return Literal{Value: p.previous().Literal}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{Expression: expr}
	}
	errorToken(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}
func (p *Parser) consume(t int, message string) {
	if p.check(t) {
		p.advance()
		return
	}
	errorToken(p.peek(), message)
}
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}
		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

func errorToken(token Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
	}
}
