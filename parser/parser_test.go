package parser

import (
	"testing"
	"fmt"
	"yokan/ast"
	"yokan/lexer"
)


// 文のテスト

func TestAssignStatement(t *testing.T) {
	input := "aaa = \"bbb\"\nccc = 1 + 2 == ddd"

	program := checkCommonTestsAndParse(t, input, 2)

	tests := []struct {
		expectedIdentifier string
		expectedExpression string
	} {
		{"aaa", `"bbb"`},
		{"ccc", "((1 + 2) == ddd)"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !checkAssignStatement(t, stmt, tt.expectedIdentifier, tt.expectedExpression) {
			return
		}
	}
}

func checkAssignStatement(t *testing.T, stmt ast.Statement, name string, expected string) bool {
	assign, ok := stmt.(*ast.Assign)
	if !ok {
		t.Errorf("stmt is not *ast.Assign. got=%T", stmt)
		return false
	}
	checkIdentifier(t, &assign.Name, name)
	if assign.Value.String() != expected {
		t.Errorf("assign.Value.String() is not '%s'. got %s", expected, assign.Value.String())
	}
	return true
}

// 式のテスト

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
			t.Fatalf("expr is not ast.InfixExpression. got %T", expr)
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
	tests := []testInString {
		{"-a * b", "((-a) * b)"},
		{"--a", "(-(-a))"},
		{"-+a", "(-(+a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + -b - c", "((a + (-b)) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * +b / c", "((a * (+b)) / c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 <= 4 != 3 >= 4", "((5 <= 4) != (3 >= 4))"},
		{"a * b + c * d", "((a * b) + (c * d))"},
		{"1 == 2", "(1 == 2)"},

		{"z == a * b + c", "(z == ((a * b) + c))"},
		{"z == a + b * c", "(z == (a + (b * c)))"},
	}
	checkExpressionsInString(t, tests)
}

func TestParenthesisExpressions(t *testing.T) {
	tests := []testInString {
		{"a * (b + c)", "(a * (b + c))"},
		{"(a + b) * c", "((a + b) * c)"},
		{"(a + ((b > c) == (d <= e))) * f", "((a + ((b > c) == (d <= e))) * f)"},
		{"a * (b != c)", "(a * (b != c))"},
	}
	checkExpressionsInString(t, tests)
}

type testInString struct {
	input string
	expected string
}

func checkExpressionsInString(t *testing.T, tests []testInString) {
	for _, tt := range tests {
		expr := checkCommonTestsAndParseExpression(t, tt.input)
		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestFunctionCalling(t *testing.T) {
	input := "add(x,2)"
	
	expr := checkCommonTestsAndParseExpression(t, input)
	
	calling, ok := expr.(*ast.FunctionCalling)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionCalling. got=%T", expr)
	}
	checkIdentifier(t, calling.Function, "add")
	args := calling.Arguments
	if len(args) != 2 {
		t.Fatalf("len(args) is not 2. got=%q", args)
	}
	checkIdentifier(t, args[0], "x")
	checkIntegerLiteral(t, args[1], 2)
}

func TestFunctionCalling2(t *testing.T) {
	input := "func(a)(b)"
	
	expr := checkCommonTestsAndParseExpression(t, input)
	
	outerCalling, ok := expr.(*ast.FunctionCalling)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionCalling. got=%T", expr)
	}
	outerArgs := outerCalling.Arguments
	if len(outerArgs) != 1 {
		t.Fatalf("len(outerArgs) is not 1. got=%q", outerArgs)
	}
	checkIdentifier(t, outerArgs[0], "b")

	innerCalling, ok := outerCalling.Function.(*ast.FunctionCalling)
	if !ok {
		t.Fatalf("innerCalling.Function is not *ast.FunctionCalling. got=%T", innerCalling.Function)
	}
	innerArgs := innerCalling.Arguments
	if len(innerArgs) != 1 {
		t.Fatalf("len(innerArgs) is not 1. got=%q", innerArgs)
	}
	checkIdentifier(t, innerArgs[0], "a")

	checkIdentifier(t, innerCalling.Function, "func")
}

func TestFunctionLiteral(t *testing.T) {
	input := "(aa,bb,){\ncc=44\nee\n}"

	expr := checkCommonTestsAndParseExpression(t, input)

	fun, ok := expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionLiteral. got=%T", expr)
	}
	args := fun.Arguments
	if len(args) != 2 {
		t.Fatalf("len(args) is not 2. got=%q", args)
	}
	checkIdentifier(t, &args[0], "aa")
	checkIdentifier(t, &args[1], "bb")
	first_stmt, ok := fun.Body[0].(*ast.Assign)
	if !ok {
		t.Fatalf("fun.Body[0] is not *ast.Assign. got=%T", fun.Body[0])
	}
	checkIdentifier(t, &first_stmt.Name, "cc")
	checkIntegerLiteral(t, first_stmt.Value, 44)
	second_stmt, ok := fun.Body[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fun.Body[1] is not *ast.ExpressionStatement. got=%T", fun.Body[1])
	}
	checkIdentifier(t, second_stmt.Expression, "ee")
}

func TestFunctionLiteralWithCalling(t *testing.T) {
	input := "(){}()"

	expr := checkCommonTestsAndParseExpression(t, input)

	call, ok := expr.(*ast.FunctionCalling)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionCalling. got=%T", expr)
	}
	fun, ok := call.Function.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("call.Function is not *ast.FunctionLiteral. got=%T", call.Function)
	}
	if len(fun.Arguments) != 0 {
		t.Fatalf("len(fun.Arguments) is not 0. got=%d", len(fun.Arguments))
	}
	if len(fun.Body) != 0 {
		t.Fatalf("len(fun.Body) is not 0. got=%d", len(fun.Body))
	}
}

func TestFunctionLiteralWithCalling2(t *testing.T) {
	input := "(a){a}(12)"

	expr := checkCommonTestsAndParseExpression(t, input)

	call, ok := expr.(*ast.FunctionCalling)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionCalling. got=%T", expr)
	}
	fun, ok := call.Function.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("call.Function is not *ast.FunctionLiteral. got=%T", call.Function)
	}
	if len(fun.Arguments) != 1 {
		t.Fatalf("len(fun.Arguments) is not 0. got=%d", len(fun.Arguments))
	}
	if len(fun.Body) != 1 {
		t.Fatalf("len(fun.Body) is not 0. got=%d", len(fun.Body))
	}
}

func TestFunctionLiteralWithCalling3(t *testing.T) {
	input := "(){ffff()}\n 123"

	exprs := checkCommonTestsAndParse(t, input, 2)
	expr, ok := exprs.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exprs.Statements[1] is not *ast.ExpressionStatement. got=%T", exprs.Statements[1])
	}
	checkIntegerLiteral(t, expr.Expression, 123)
}

// リテラルのテスト

func TestIntegerLiteralExpression(t *testing.T) {
	input := "11"

	expr := checkCommonTestsAndParseExpression(t, input)

	checkIntegerLiteral(t, expr, 11)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"aa\n\t\"a"`

	expr := checkCommonTestsAndParseExpression(t, input)

	checkStringLiteral(t, expr, "aa\n\t\"a")	
}

func TestArrayLiteralExperession(t *testing.T) {
	input := `[12, "bb", [33, [], ]]`

	expr := checkCommonTestsAndParseExpression(t, input)

	array, ok := expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("array is not *ast.ArrayLiteral. got=%T", array)
	}
	if len(array.Value) != 3 {
		t.Fatalf("len(array.Value) is not 3. got=%d", len(array.Value))
	}
	checkIntegerLiteral(t, array.Value[0], 12)
	checkStringLiteral(t, array.Value[1], "bb")
	
	innerArray, ok := array.Value[2].(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("innerArray is not *ast.ArrayLiteral. got=%T", innerArray)
	}
	if len(innerArray.Value) != 2 {
		t.Fatalf("len(innerArray.Value) is not 2. got=%d", len(array.Value))
	}
	checkIntegerLiteral(t, innerArray.Value[0], 33)

	innerInnerArray, ok := innerArray.Value[1].(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("innerInnerArray is not *ast.ArrayLiteral. got=%T", innerInnerArray)
	}
	if len(innerInnerArray.Value) != 0 {
		t.Fatalf("len(innerArray.Value) is not 0. got=%d", len(innerInnerArray.Value))
	}
}

// リテラルと識別子のチェック

func checkIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp is not *ast.IntegerLiteral. got=%T", exp)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value is not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral is not %d. got %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func checkStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	literal, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.StringLiteral. got=%T", exp)
		return false
	}
	if literal.Value != value {
		t.Fatalf("Literal.Value is not %s. got=%s", value, literal.Value)
		return false
	}
	if literal.TokenLiteral() != value {
		t.Fatalf("literal.TokenLiteral is not %s. got %s",
			value, literal.TokenLiteral())
		return false
	}
	return true
}

func checkIdentifier(t *testing.T, exp ast.Expression, name string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Name != name {
		t.Fatalf("ident.Name is not '%s'. got='%s'", name, ident.Name)
		return false
	}
	if ident.TokenLiteral() != name {
		t.Fatalf("ident.TokenLiteral() is not '%s'. got='%s'", name, ident.TokenLiteral())
		return false
	}
	return true
}

// 式や文のテストに使う共通部分

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
		t.Fatalf("program.Statements does is not contain %d statements. got=%d, %q",
			neededStmt, len(program.Statements), program.Statements)
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
