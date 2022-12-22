package parser

import (
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
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	return nil
}
