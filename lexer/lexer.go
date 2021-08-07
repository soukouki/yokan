package lexer

import (
	"yokan/token"
)

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch==' ' || l.ch=='\t' || l.ch=='\r' {
		l.readChar()
	}
}

func (l *Lexer) skipNewlines() {
	for l.ch=='\n' {
		l.readChar()
	}
}

func (l *Lexer) skipLines() {
	for l.ch!='\n' {
		l.readChar()
	}

}

func (l *Lexer) nextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpaces()
	
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.skipLines()
			tok = newToken(token.NEWLINE, '\n')
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
		l.skipNewlines()
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = "EOF"
	default:
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readDigits()
			return tok
		} else if isLetter(l.ch) {
			tok.Type = token.IDENT
			tok.Literal = l.readIdentifier()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func isDigit(ch byte) bool {
	return include('0', '9', ch)
}

func isLetter(ch byte) bool {
	return include('a', 'z', ch) || include('A', 'Z', ch) || ch == '_'
}

func include(start byte, end byte, ch byte) bool {
	return start <= ch && ch <= end
}

func (l *Lexer) readDigits() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
