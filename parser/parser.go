package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

// we need to look at the curToken, which is the current token under
// examination, to decide what to do next, and we also need peekToken for this decision if curToken
// doesn’t give us enough information
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	//error
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	//init curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}
	return nil
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	//let xx = yyyyy
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}
	//expect a identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	//expect '='
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO
	//cope with expression

	//skip expression until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	//return <expression>

	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	//TODO:deal with ReturnValue

	//skip expression until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
