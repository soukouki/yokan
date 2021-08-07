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

	testTokens(t, input, expected)
}

func TestSpaces(t *testing.T) {
	input := " \t \n "
	expected := []TypeAndLiteral {
		{token.NEWLINE, "\n"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestComment(t *testing.T) {
	input := "+//aaa\n-"
	expected := []TypeAndLiteral {
		{token.PLUS, "+"},
		{token.NEWLINE, "\n"},
		{token.MINUS, "-"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestContiguousComments(t *testing.T) {
	input := "+//aaa\n//bbb\n-"
	expected := []TypeAndLiteral {
		{token.PLUS, "+"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"}, // トークナイザで処理するのは諦める
		{token.MINUS, "-"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestInteger(t *testing.T) {
	input := "1+234*567 89"
	expected := []TypeAndLiteral {
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "234"},
		{token.STAR, "*"},
		{token.INT, "567"},
		{token.INT, "89"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestIdentifier(t *testing.T) {
	input := "a+bbb*CcC ddddd _eE_123e"
	expected := []TypeAndLiteral{
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "bbb"},
		{token.STAR, "*"},
		{token.IDENT, "CcC"},
		{token.IDENT, "ddddd"},
		{token.IDENT, "_eE_123e"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func testTokens(t *testing.T, input string, expected []TypeAndLiteral) {
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

