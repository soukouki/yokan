package ast

import (
	"yokan/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() { }
func (es *ExpressionStatement) TokenLiteral() string {
	return ""
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Assign struct {
	Name *Identifier
	Value Expression
}

func (as *Assign) expressionNode() { }
func (as *Assign) TokenLiteral() string {
	return ""
}

func (as *Assign) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String()+" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString("\n")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() { }
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() { }
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}