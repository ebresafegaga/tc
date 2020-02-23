package parser

import (
	"fmt"
	"strconv"

	"github.com/ebresafegaga/tc/ast"
	"github.com/ebresafegaga/tc/lexer"
	"github.com/ebresafegaga/tc/token"
)

// Parser for the source text
type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New creates a new parser
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// Read two, so that cur and peek are both set!
	p.nextToken()
	p.nextToken()

	return p
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

// Errors returns a slice of error messages
func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

// ParseProgram parses a program into an ast
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.curToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}
	return program
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// 	Looping over the expression for now
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("Expected next token to be %s, but got %s intead", t, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) curTokenIs(t token.Type) bool {
	return parser.curToken.Type == t
}

func (parser *Parser) peekTokenIs(t token.Type) bool {
	return parser.peekToken.Type == t
}

func (parser *Parser) expectPeek(t token.Type) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	}

	// we probably have a syntax error
	parser.peekError(t)
	return false
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()

	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}
	return stmt
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.curToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.curToken.Type]
	if prefix == nil {
		return nil
	}

	left := prefix()

	return left
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
}

// Pratt Parsing

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}
