package ast

import (
	"bytes"
	"yokan/utility"
	"yokan/token"
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


// プログラム(全体)

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


// 式だけの文

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


// 代入(文)

type Assign struct {
	Name Identifier
	Value Expression
}

func (as *Assign) statementNode() { }
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


// 前置演算子

type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}

func (p *PrefixExpression) expressionNode() { }
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}


// 中置演算子

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (i *InfixExpression) expressionNode() { }
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" ")
	out.WriteString(i.Operator)
	out.WriteString(" ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}


// 関数呼び出し

type FunctionCalling struct {
	Token token.Token
	Function Expression
	Arguments []Expression
}

func (fc *FunctionCalling) expressionNode() { }
func (fc *FunctionCalling) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc *FunctionCalling) String() string {
	var out bytes.Buffer
	out.WriteString(fc.Function.String())
	out.WriteString("(")
	len := len(fc.Arguments)
	for i, e := range fc.Arguments {
		out.WriteString(e.String())
		if i != len {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}


// 配列

type ArrayLiteral struct {
	Token token.Token
	Value []Expression
}

func (a *ArrayLiteral) expressionNode() { }
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	len := len(a.Value)
	for i, e := range a.Value {
		out.WriteString(e.String())
		if i != len {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")
	return out.String()
}


// 関数リテラル

type FunctionLiteral struct {
	Token token.Token
	Arguments []Identifier
	Body []Statement
}

func (f *FunctionLiteral) expressionNode() { }
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	identLen := len(f.Arguments)
	for i, ident := range f.Arguments {
		out.WriteString(ident.String())
		if i != identLen {
			out.WriteString(", ")
		}
	}
	out.WriteString(")\n{\n")
	for _, stmt := range f.Body {
		out.WriteString("\t")
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}


// 識別子

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


// 整数リテラル

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


// 文字列リテラル

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() { }
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return utility.Quote(sl.Value)
}
