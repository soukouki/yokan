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

func TestPrefixExpressions(t *testing.T) {
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

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not Integer. got=%T", obj)
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
