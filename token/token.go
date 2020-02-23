package token

// Type represents
type Type string

// Token represents
type Token struct {
	Type    Type
	Literal string
}

// SyntaxKind represents
type SyntaxKind string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	EQ    = "=="
	NOTEQ = "!="

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ESLE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ESLE,
	"return": RETURN,
}

// LookupIdent checks
func LookupIdent(ident string) Type {
	if res, ok := keywords[ident]; ok {
		return res
	}
	return IDENT
}
