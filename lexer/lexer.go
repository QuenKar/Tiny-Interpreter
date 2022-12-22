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
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) || l.readPosition < 0 {
		// sets l.ch to 0, which is the ASCII code for the "NUL" character
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

//获取下一个pos的char
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) || l.readPosition < 0 {
		// return the ASCII code for the "NUL" character
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token
	//skip the white space
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			c := l.ch
			l.readChar()
			t = token.Token{Type: token.EQ, Literal: string(c) + string(l.ch)}
		} else {
			t = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '!':
		if l.peekChar() == '=' {
			c := l.ch
			l.readChar()
			t = token.Token{Type: token.NOT_EQ, Literal: string(c) + string(l.ch)}
		} else {
			t = newToken(token.BANG, l.ch)
		}
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.MINUS, l.ch)
	case '*':
		t = newToken(token.ASTERISK, l.ch)
	case '/':
		t = newToken(token.SLASH, l.ch)
	case '<':
		t = newToken(token.LT, l.ch)
	case '>':
		t = newToken(token.GT, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Type = token.CheckIdent(t.Literal)
			return t

		} else if isDigit(l.ch) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		} else {
			t = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return t
}
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func isLetter(ch byte) bool {
	//字母和下划线
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
