package parser

import (
	"testing"
	"fmt"
	"yokan/ast"
	"yokan/lexer"
)

func TestAssignStatement(t *testing.T) {
	input := "aaa = 123\nbbb = \"ccc\""

	program := checkCommonTestsAndParse(t, input, 2)

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

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input string
		left int64
		operator string
		right int64
	} {
		{"1+2", 1, "+", 2},
		{"3-4", 3, "-", 4},
		{"5*6", 5, "*", 6},
		{"7/8", 7, "/", 8},
		{"9==10", 9, "==", 10},
		{"11!=12", 11, "!=", 12},
		{"13<14", 13, "<", 14},
		{"15<=16", 15, "<=", 16},
		{"17>18", 17, ">", 18},
		{"19>=20", 19, ">=", 20},
	}

	for _, tt := range infixTests {
		expr := checkCommonTestsAndParseExpression(t, tt.input)
		ie, ok := expr.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got %T", expr)
		}
		if !checkIntegerLiteral(t, ie.Left, tt.left) {
			return
		}
		if ie.Operator != tt.operator {
			t.Fatalf("ie.Operator is not '%s'. got=%s", tt.operator, ie.Operator)
		}
		if !checkIntegerLiteral(t, ie.Right, tt.right) {
			return
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		integerValue int64
	} {
		{"+12", "+", 12},
		{"-34", "-", 34},
	}

	for _, tt := range prefixTests {
		expr := checkCommonTestsAndParseExpression(t, tt.input)
		pe, ok := expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", expr)
		}
		if pe.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, pe.Operator)
		}
		if !checkIntegerLiteral(t, pe.Right, tt.integerValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{"-a * b", "((-a) * b)"},
		{"--a", "(-(-a))"},
		{"-+a", "(-(+a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + -b - c", "((a + (-b)) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * +b / c", "((a * (+b)) / c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		//{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		//{"5 <= 4 != 3 >= 4", "((5 <= 4) != (3 <= 4))"},
		{"1 == 2", "(1 == 2)"},

		//{"z == a * b + c", "(z == ((a * b) + c))"},
		//{"z == a + b * c", "(z == (a + (b * c)))"},
	}

	for _, tt := range tests {
		expr := checkCommonTestsAndParseExpression(t, tt.input)
		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "11"

	expr := checkCommonTestsAndParseExpression(t, input)

	checkIntegerLiteral(t, expr, 11)
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"aa\n\t\"a"`

	expr := checkCommonTestsAndParseExpression(t, input)

	literal, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T, stmt.Expression",
			expr)
	}
	if literal.Value != "aa\n\t\"a" {
		t.Fatalf("Literal.Value not %s. got=%s", "aa", literal.Value)
	}
	if literal.TokenLiteral() != "aa\n\t\"a" {
		t.Fatalf("literal.TokenLiteral not %s. got %s",
			"aa\n\t\"a", literal.TokenLiteral())
	}
}

func checkCommonTestsAndParseExpression(t *testing.T, input string) ast.Expression {
	program := checkCommonTestsAndParse(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}

	return stmt.Expression
}

func checkCommonTestsAndParse(t *testing.T, input string, neededStmt int) *ast.Program {
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
