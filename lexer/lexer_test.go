package lexer

import (
	"testing"
	"yokan/token"
)

type TypeAndLiteral struct {
	Type token.TokenType
	Literal string
}

func TestOneCharacterKeywords(t *testing.T) {
	input := `=+-*/,(){}<>[]`

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
		{token.LBRACK, "["},
		{token.RBRACK, "]"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestTwoCharacterKeywords(t *testing.T) {
	input := "== != <= >="
	expected := []TypeAndLiteral {
		{token.EQ, "=="},
		{token.NOTEQ, "!="},
		{token.LTEQ, "<="},
		{token.GTEQ, ">="},
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
	input := "a+bbb*CcC ddddd _eE_123e\nfff"
	expected := []TypeAndLiteral {
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "bbb"},
		{token.STAR, "*"},
		{token.IDENT, "CcC"},
		{token.IDENT, "ddddd"},
		{token.IDENT, "_eE_123e"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "fff"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestString(t *testing.T) {
	input := "\"abc\" \"\" \"\\\"\" \"\\n\\t\" \"\n\""
	expected := []TypeAndLiteral {
		{token.STRING, "abc"},
		{token.STRING, ""},
		{token.STRING, "\""},
		{token.STRING, "\n\t"},
		{token.STRING, "\n"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestKeyWords(t *testing.T) {
	input := "func return while if else"
	expected := []TypeAndLiteral {
		{token.FUNC, "func"},
		{token.RETURN, "return"},
		{token.WHILE, "while"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.EOF, "EOF"},
	}
	testTokens(t, input, expected)
}

func TestComplexSource(t *testing.T) {
	input := `
		a = 123 + 456 // コメント
		a = a+789
		str = "abab\n" // 特殊文字はとりあえず改行だけ
		array = [1, 2]

		func add(x, y) {
		  x + y
		}

		i = 0
		while(i<4) {
		  val = array[i]
		  if(val <= 2) {
		    print(val)
		  } else {
		    print(val*10)
		  }
		}
	`
	expected := []TypeAndLiteral {
		{token.NEWLINE, "\n"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "123"},
		{token.PLUS, "+"},
		{token.INT, "456"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.INT, "789"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "str"},
		{token.ASSIGN, "="},
		{token.STRING, "abab\n"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "array"},
		{token.ASSIGN, "="},
		{token.LBRACK, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACK, "]"},
		{token.NEWLINE, "\n"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "val"},
		{token.ASSIGN, "="},
		{token.IDENT, "array"},
		{token.LBRACK, "["},
		{token.IDENT, "i"},
		{token.RBRACK, "]"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "val"},
		{token.LTEQ, "<="},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "val"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "val"},
		{token.STAR, "*"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
	}
	testTokens(t, input, expected)
}

func testTokens(t *testing.T, input string, expected []TypeAndLiteral) {
	l := New(input)

	for i, expected := range expected {
		tok := l.NextToken()
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

