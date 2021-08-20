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

// TODO: 比較

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not *Integer. got=%T", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
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
