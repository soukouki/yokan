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
	assignStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}
	if assignStmt.Name.Value != name {
		t.Errorf("assignStmt.Name.Value not '%s'. got=%s", name, assignStmt.Name.Value)
		return false
	}
	if assignStmt.Name.TokenLiteral() != name {
		t.Errorf("assignStmt.Name.TokenLiteral() not '%s'. got=%s", name, assignStmt.Name.TokenLiteral())
		return false
	}
	return true
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
