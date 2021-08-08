package parser

import (
	"fmt"
	"strconv"
	"yokan/ast"
	"yokan/lexer"
	"yokan/token"
)

type Parser struct {
	l *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string {},
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.appendError(msg)
}

func (p *Parser) appendError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() *ast.ExpressionStatement {
	var expr ast.Expression
	switch p.curToken.Type {
	case token.NEWLINE:
		return nil
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			name := p.parseIdentifier()
			expr = p.parseAssign(*name)
		} else {
			expr = p.parseInfixExpression()
		}
	default:
		expr = p.parseInfixExpression()
	}
	fmt.Print(expr, "\n")
	return &ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parseAssign(name ast.Identifier) *ast.Assign {
	assign := &ast.Assign{Token: name.Token, Name: &name}
	p.nextToken()
	assign.Value = p.parseInfixExpression()
	return assign
}

func (p *Parser) parseInfixExpression() ast.Expression {
	left := p.parsePrefixExpression()
	switch p.peekToken.Type {
	case
		token.PLUS, token.MINUS, token.STAR, token.SLASH,
		token.EQ, token.NOTEQ,
		token.LT, token.LTEQ, token.GT, token.GTEQ:
	default:
		return left
	}
	p.nextToken()
	ie := &ast.InfixExpression{Token: p.curToken}
	ie.Left = left
	ie.Operator = p.curToken.Literal
	p.nextToken()
	ie.Right = p.parseInfixExpression()
	return ie
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	switch p.curToken.Type {
	case token.PLUS, token.MINUS:
	default:
		lit := p.parseLiteralAndIdentify()
		if lit != nil {
			return lit
		}
	}
	pe := &ast.PrefixExpression{Token: p.curToken}
	pe.Operator = p.curToken.Literal
	p.nextToken()
	pe.Right = p.parsePrefixExpression()
	return pe
}

func (p *Parser) parseLiteralAndIdentify() ast.Expression {
	switch p.curToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.STRING:
		return p.parseStringLiteral()
	case token.IDENT:
		return p.parseIdentifier()
	default:
		return nil
	}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.appendError(msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
