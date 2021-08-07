package lexer

import (
	"testing"
	"yokan/token"
)

type TypeAndLiteral struct {
	Type token.TokenType
	Literal string
}

func TestOneCharacters(t *testing.T) {
	input := `=+-*/,(){}<>`

	expected := []TypeAndLiteral {
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.COMMA, ","},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.EOF, "EOF"},
	}

	_TestTokens(t, input, expected)
}

func TestSpaces(t *testing.T) {
	input := " \t \n "
	expected := []TypeAndLiteral {
		{token.NEWLINE, "\n"},
		{token.EOF, "EOF"},
	}
	_TestTokens(t, input, expected)
}

func TestComment(t *testing.T) {
	input := "+//aaa\n-"
	expected := []TypeAndLiteral {
		{token.PLUS, "+"},
		{token.NEWLINE, "\n"},
		{token.MINUS, "-"},
		{token.EOF, "EOF"},
	}
	_TestTokens(t, input, expected)
}

func _TestTokens(t *testing.T, input string, expected []TypeAndLiteral) {
	l := New(input)

	for i, expected := range expected {
		tok := l.nextToken()
		if tok.Type != expected.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expected.Type, tok.Type)
		}
		if tok.Literal != expected.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expected.Literal, tok.Literal)
		}
	}
}

