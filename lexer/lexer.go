package lexer

import "interpreter/token"

//readPosition always points to the “next” character in the input.
//position points to the character in the input that corresponds to the ch byte.
type Lexer struct {
	input        string
	position     int  //current position in input (points to current char)
	readPosition int  //current reading position in input
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input) || l.readPosition < 0 {
		// sets l.ch to 0, which is the ASCII code for the "NUL" character
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	switch l.ch {
	case '=':
		t = newToken(token.ASSIGN, l.ch)
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	}
	l.ReadChar()
	return t
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
