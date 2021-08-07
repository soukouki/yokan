package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string	
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 リテラル
	IDENT = "IDENT" // add, foobar, x, ...
	INT   = "INT"

	// 演算子
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"

	EQ     = "=="
	NOTEQ  = "!="
	LT     = "<"
	LTEQ   = "<="
	GT     = ">"
	GTEQ   = ">="


	// デミリタ
	NEWLINE = "\n"
	COMMA   = ","
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"

	// キーワード
	FUNC   = "FUNCTION"
	RETURN = "RETURN"
	WHILE  = "WHILE"
	IF     = "IF"
	ELSE   = "ELSE"
)