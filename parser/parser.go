package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

//The blank identifier _ takes the zero value and
//the following constants get assigned the values 1 to 7.
//Which numbers we use doesn’t matter, but the order matters.
// it means precedence. for example '+' < '*', and function call
//has the heightest precedence.
const (
	_           int = iota
	LOWEST          //default lowest precedence
	EQUALS          //== or !=
	LESSGREATER     //> or <
	SUM             //+
	PRODUCT         //* or /
	PREFIX          //!x or -x...
	CALL            //myfunc(x)
)

// the parsing functions are
// called to parse the appropriate expression and return an AST node that represents it
type (
	prefixParseFn func() ast.Expression               //no need of argument
	infixParseFn  func(ast.Expression) ast.Expression //This argument is “left side” of the infix operator that’s being parsed.
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// we need to look at the curToken, which is the current token under
// examination, to decide what to do next, and we also need peekToken for this decision if curToken
// doesn’t give us enough information
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	//error
	errors []string

	//In order for our parser to get the correct prefixParseFn or infixParseFn for the current token type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	//init curToken and peekToken
	p.nextToken()
	p.nextToken()
	//for prefix
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	//Identifiers
	p.registerPrefix(token.IDENT, p.parseIdentfier)
	//Integer
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// !
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	// -
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	//true and false
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// ( 2 + 3 )
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	//for infix
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
	default:
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

//How to understand "precedence" parameter ???
//right-binding power: the higher it is, the more tokens/operators/operands
//to the right of the current expressions (the future peek tokens) can we “bind” to it
//think about Expression:  1 + 2 * 3
func (p *Parser) parseExpression(precedence int) ast.Expression {
	//deal with prefix
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	//deal with infix
	//precedence < p.peekPrecedence():This condition checks if the left-binding power of
	//the next operator/token is higher than our current right-binding power. I
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

//returns a *ast.Identifier with the current token in the Token field and the literal value of the token in Value
func (p *Parser) parseIdentfier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: p.curToken,
	}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	literal.Value = value
	return literal
}

/*
When parsePrefixExpression is called, p.curToken is either of type
token.BANG or token.MINUS,But in order to correctly parse a prefix
expression like -5 more than one token has to be “consumed”. So
after using p.curToken to build a *ast.PrefixExpression node, the
method advances the tokens and calls parseExpression again.
*/
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	curP := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(curP)

	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", t))
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

/*
Think about ( 2 + 3 )
when encounter ')', p.parseExpression function will return
because ')' precedence is LOWEST
*/
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}
