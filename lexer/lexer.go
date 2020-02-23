package lexer

import (
	"github.com/ebresafegaga/tc/token"
)

// Lexer represents the basic state of the lexer.
type Lexer struct {
	input        string
	position     int  // Current ch position
	readPosition int  // Position after current ch, for reading
	ch           byte // Current ch under examination
}

// New returns a new new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.ch = 0
	} else {
		lexer.ch = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition++
}

// NextToken get the next lexeme in the source text
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.eatWhitespaces()

	switch lexer.ch {
	case '=':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(lexer.ch)}
		} else {
			tok = newToken(token.ASSIGN, lexer.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case '-':
		tok = newToken(token.MINUS, lexer.ch)
	case '!':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			tok = token.Token{Type: token.NOTEQ, Literal: string(ch) + string(lexer.ch)}
		} else {
			tok = newToken(token.BANG, lexer.ch)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.ch)
	case '*':
		tok = newToken(token.ASTERISK, lexer.ch)
	case '<':
		tok = newToken(token.LT, lexer.ch)
	case '>':
		tok = newToken(token.GT, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.ch) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}
	lexer.readChar()

	return tok
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (lexer *Lexer) readIdentifier() string {
	start := lexer.position
	for isLetter(lexer.ch) {
		lexer.readChar()
	}
	end := lexer.position
	return lexer.input[start:end]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) eatWhitespaces() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	}
	return lexer.input[lexer.readPosition]
}
