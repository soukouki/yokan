package evaluator

import (
	"testing"
	"yokan/lexer"
	"yokan/parser"
	"yokan/object"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaled := testEval(tt.input)
		testIntegerObject(t, evaled, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	evaled := testEval(`"a\n\t\"b"`)
	str, ok := evaled.(*object.String)
	if !ok {
		t.Errorf("evaled is not object.String. got=%T", evaled)
	}
	if str.Value != "a\n\t\"b" {
		t.Errorf(`str.Value is not "a\n\t\"b". got='%s'`, str)
	}
}

func TestEvalPrefixExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"+12", 12},
		{"-34", -34},
		{"-+-+56", 56},
	}

	for _, tt := range tests {
		evaled := testEval(tt.input)
		testIntegerObject(t, evaled, tt.expected)
	}
}

func TestPrefixExpressionsTypeError(t *testing.T) {
	evaled := testEval(`+""`)
	err, ok := evaled.(*object.TypeMisMatchError)
	if !ok {
		t.Errorf("evaled is not *object.TypeMisMatchError. got=%T", evaled)
	}
	if err.Type() != "ERROR" {
		t.Errorf("err.Type() is not 'ERROR. got='%s'", err.Type())
	}
}

func TestEvalInfixIntegerExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"1+2", 3},
		{"3-4", -1},
		{"5*6", 30},
		{"7/8", 0},
	}
	for _, tt := range tests {
		evaled := testEval(tt.input)
		testIntegerObject(t, evaled, tt.expected)
	}
}

func TestZeroDivisionError(t *testing.T) {
	evaled := testEval("1/0")
	_, ok := evaled.(*object.OtherError)
	if !ok {
		t.Errorf("evaled is not *object.OtherError. got=%T", evaled)
	}
	// メッセージのチェックはとりあえずしない
}

func TestEvalInfixComparingExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	} {
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 1", false},
		{"1 != 2", true},
		{`"a" == "a"`, true},
		{`"a" == "b"`, false},
		{`"a" != "a"`, false},
		{`"a" != "b"`, true},
		{"(1==1) == (1==1)", true},
		{"(1==1) == (1!=1)", false},
		{"(1==1) != (1==1)", false},
		{"(1==1) != (1!=1)", true},
		{`1 == "a"`, false},
		{`1 != "a"`, true},
		{"1 == (1==1)", false},
		{"1 != (1==1)", true},
		{"1 < 2", true},
		{"4 < 3", false},
		{"1 < 1", false},
		{"1 <= 2", true},
		{"4 <= 3", false},
		{"1 <= 1", true},
		{"1 > 2", false},
		{"4 > 3", true},
		{"1 > 1", false},
		{"1 >= 2", false},
		{"4 >= 3", true},
		{"1 >= 1", true},
	}
	for _, tt := range tests {
		evaled := testEval(tt.input)
		testBooleanObject(t, evaled, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not *object.Integer. got=%T(%s)", obj, obj.Inspect())
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not object.Boolean. got=%T(%s)", obj, obj.Inspect())
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	return Eval(prog)
}
