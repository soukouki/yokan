package parser

import (
	"testing"
	"yokan/ast"
	"yokan/lexer"
)

func TestAssignStatement(t *testing.T) {
	input := "aaa = 123\nbbb = \"ccc\""

	program := checkCommonAndGenerateProgram(t, input, 2)

	tests := []struct {
		expectedIdentifier string
	} {
		{"aaa"},
		{"bbb"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !checkAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	exprStmt, ok := s.(*ast.ExpressionStatement)
	assign := exprStmt.Expression.(*ast.Assign)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}
	if assign.Name.Value != name {
		t.Errorf("assign.Name.Value not '%s'. got=%s", name, assign.Name.Value)
		return false
	}
	if assign.Name.TokenLiteral() != name {
		t.Errorf("assign.Name.TokenLiteral() not '%s'. got=%s", name, assign.Name.TokenLiteral())
		return false
	}
	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "11"

	program := checkCommonAndGenerateProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T, stmt.Expression",
			program.Statements[0])
	}
	if literal.Value != 11 {
		t.Fatalf("literal.Value not %d. got=%d", 11, literal.Value)
	}
	if literal.TokenLiteral() != "11" {
		t.Fatalf("literal.TokenLiteral not %s. got %s",
			"11", literal.TokenLiteral())
	}
}

func checkCommonAndGenerateProgram(t *testing.T, input string, neededStmt int) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	checkParseErrors(t, p)

	if len(program.Statements) != neededStmt {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			neededStmt, len(program.Statements))
	}

	return program
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
