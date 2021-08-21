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
	peek2Token token.Token
	peek3Token token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string {},
	}
	p.nextToken()
	p.nextToken()
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead",
		t, p.peekToken.Type)
	p.appendError(msg)
}

func (p *Parser) appendError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.peek2Token
	p.peek2Token = p.peek3Token
	p.peek3Token = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = p.parseStatements()
	return program
}

func (p *Parser) parseStatements() []ast.Statement {
	var list []ast.Statement
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			list = append(list, stmt)
		}
		p.nextToken()
	}
	return list
}

func (p *Parser) parseStatement() ast.Statement {
	var expr ast.Expression
	switch p.curToken.Type {
	case token.NEWLINE:
		return nil
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssign()
		} else {
			expr = p.parseExpression()
		}
	default:
		expr = p.parseExpression()
	}
	return &ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parseAssign() *ast.Assign {
	assign := &ast.Assign{Name: *p.parseIdentifier()}
	p.nextToken()
	p.nextToken()
	assign.Value = p.parseExpression()
	return assign
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseEqExpression()
}

func (p *Parser) parseEqExpression() ast.Expression {
	expr := p.parseLTGTExpression()
	for p.peekTokenIs(token.EQ) || p.peekTokenIs(token.NOTEQ) {
		p.nextToken()
		newExpr := &ast.InfixExpression{
			Token: p.curToken,
			Left: expr,
			Operator: p.curToken.Literal,
		}
		p.nextToken()
		newExpr.Right = p.parseLTGTExpression()
		expr = newExpr
	}
	return expr
}

func (p *Parser) parseLTGTExpression() ast.Expression {
	expr := p.parseAddSubExpression()
	for (
		p.peekTokenIs(token.LT) || p.peekTokenIs(token.LTEQ) || 
		p.peekTokenIs(token.GT) || p.peekTokenIs(token.GTEQ) ) {
		p.nextToken()
		newExpr := &ast.InfixExpression{
			Token: p.curToken,
			Left: expr,
			Operator: p.curToken.Literal,
		}
		p.nextToken()
		newExpr.Right = p.parseAddSubExpression()
		expr = newExpr
	}
	return expr
}

func (p *Parser) parseAddSubExpression() ast.Expression {
	expr := p.parseMulDivExpression()
	for p.peekTokenIs(token.PLUS) || p.peekTokenIs(token.MINUS) {
		p.nextToken()
		newExpr := &ast.InfixExpression{
			Token: p.curToken,
			Left: expr,
			Operator: p.curToken.Literal,
		}
		p.nextToken()
		newExpr.Right = p.parseMulDivExpression()
		expr = newExpr
	}
	return expr
}

func (p *Parser) parseMulDivExpression() ast.Expression {
	expr := p.parsePrefixExpression()
	for p.peekTokenIs(token.STAR) || p.peekTokenIs(token.SLASH) {
		p.nextToken()
		newExpr := &ast.InfixExpression{
			Token: p.curToken,
			Left: expr,
			Operator: p.curToken.Literal,
		}
		p.nextToken()
		newExpr.Right = p.parsePrefixExpression()
		expr = newExpr
	}
	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if !( p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) ) {
		return p.parseFunctionCalling()
	}
	pe := &ast.PrefixExpression{Token: p.curToken}
	pe.Operator = p.curToken.Literal
	p.nextToken()
	pe.Right = p.parsePrefixExpression()
	return pe
}

func (p *Parser) parseFunctionCalling() ast.Expression {
	expr := p.parseFunctionLiteral()
	for p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		fc := &ast.FunctionCalling{
			Token: p.curToken,
			Function: expr,
		}
		p.nextToken()
		fc.Arguments = p.parseCommaSeparatedExpressions(token.RPAREN)
		expr = fc
	}
	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	// (){ ... }
	// (a){ ... }
	//  (a, b){ ... } のようなコードを想定(カンマは式の中には出てこないのと、関数リテラルのカッコ内は識別子だけなのを用いる)
	if !(
		p.curTokenIs(token.LPAREN) && p.peekTokenIs(token.RPAREN) && p.peek2TokenIs(token.LBRACE) ||
		p.curTokenIs(token.LPAREN) && p.peekTokenIs(token.IDENT) && p.peek2TokenIs(token.RPAREN) && p.peek3TokenIs(token.LBRACE) ||
		p.curTokenIs(token.LPAREN) && p.peekTokenIs(token.IDENT) && p.peek2TokenIs(token.COMMA) ) {
		return p.parseParenthesisExpression()
	}

	token := p.curToken
	p.nextToken()
	args := p.parseCommaSeparatedIdentifiers()
	p.nextToken()
	p.nextToken()
	stmts := p.parseStatements()
	return &ast.FunctionLiteral{Token: token, Arguments: args, Body: stmts}
}

func (p *Parser) parseParenthesisExpression() ast.Expression {
	if !p.curTokenIs(token.LPAREN) {
		return p.parseLiteralAndIdentify()
	}

	p.nextToken()
	expr := p.parseExpression()
	p.expectPeek(token.RPAREN)
	return expr
}

func (p *Parser) parseLiteralAndIdentify() ast.Expression {
	switch p.curToken.Type {
	case token.LBRACK:
		return p.parseArrayLiteral()
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

func (p *Parser) parseArrayLiteral() *ast.ArrayLiteral {
	tok := p.curToken
	p.nextToken()
	list := p.parseCommaSeparatedExpressions(token.RBRACK)
	return &ast.ArrayLiteral{Token: tok, Value: list}
}

func (p *Parser) parseCommaSeparatedIdentifiers() []ast.Identifier {
	var list []ast.Identifier
	if !p.curTokenIs(token.IDENT){
		empty := []ast.Identifier { }
		return empty
	}
	for p.curTokenIs(token.IDENT) {
		ident := p.parseIdentifier()
		if ident == nil {
			msg := fmt.Sprintf("could is not parse %q as identifier", p.curToken.Literal)
			p.appendError(msg)
		}
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
		list = append(list, *ident)
	}
	return list
}

func (p *Parser) parseCommaSeparatedExpressions(endToken token.TokenType) []ast.Expression {
	var list []ast.Expression
	if p.curTokenIs(endToken) {
		p.nextToken()
		empty := []ast.Expression { }
		return empty
	}
	for {
		expr := p.parseExpression()
		p.nextToken()
		if expr==nil {
			break
		}
		list = append(list, expr)
		// , ...
		// ]
		// , ]
		if p.curTokenIs(endToken) {
			break
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
			if p.curTokenIs(endToken) {
				p.nextToken()
				break
			}
		}
	}
	return list
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could is not parse %q as integer", p.curToken.Literal)
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
		Name: p.curToken.Literal,
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peek2TokenIs(t token.TokenType) bool {
	return p.peek2Token.Type == t
}

func (p *Parser) peek3TokenIs(t token.TokenType) bool {
	return p.peek3Token.Type == t
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
