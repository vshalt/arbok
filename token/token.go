package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ASSIGN       = "="
	BANG         = "!"
	PLUS         = "+"
	MINUS        = "-"
	ASTERISK     = "*"
	SLASH        = "/"
	GREATER_THAN = ">"
	LESS_THAN    = "<"
	COMMA        = ","
	SEMICOLON    = ";"
	LPAREN       = "("
	RPAREN       = ")"
	LBRACE       = "{"
	RBRACE       = "}"

	EQUAL     = "=="
	NOT_EQUAL = "!="

	LET        = "LET"
	FUNCTION   = "FUNCTION"
	IF         = "IF"
	ELSE       = "ELSE"
	RETURN     = "RETURN"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	INTEGER    = "INTEGER"
	IDENTIFIER = "IDENTIFIER"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"func":   FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
