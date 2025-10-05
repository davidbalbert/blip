package lex

import "fmt"

type TokenType uint8

const (
	TokInvalid TokenType = iota
	TokEOF
	TokInt
	TokAdd
	TokSub
	TokMul
	TokDiv
	TokLParen
	TokRParen
)

func (t TokenType) String() string {
	switch t {
	case TokInvalid:
		return "Invalid"
	case TokEOF:
		return "EOF"
	case TokInt:
		return "Int"
	case TokAdd:
		return "Add"
	case TokSub:
		return "Sub"
	case TokMul:
		return "Mul"
	case TokDiv:
		return "Div"
	case TokLParen:
		return "LParen"
	case TokRParen:
		return "RParen"
	default:
		return "Unknown"
	}
}

type Token uint32

func NewToken(TokenType TokenType, pos uint32) Token {
	pos = pos & 0xFFFFFF
	return Token((uint32(TokenType) << 24) | pos)
}

func (t Token) Type() TokenType {
	return TokenType(t >> 24)
}

func (t Token) Pos() uint32 {
	return uint32(t) & 0xFFFFFF
}

func (t Token) String() string {
	return fmt.Sprintf("{%s Pos=%d}", t.Type(), t.Pos())
}

type Lexer struct {
	text []byte
	pos  uint32
}

func NewLexer(text []byte) *Lexer {
	return &Lexer{
		text: text,
		pos:  0,
	}
}

func (l *Lexer) NextToken() Token {
	for l.pos < uint32(len(l.text)) && isWhitespace(l.text[l.pos]) {
		l.pos++
	}

	if l.pos >= uint32(len(l.text)) {
		return NewToken(TokEOF, l.pos)
	}

	start := l.pos
	ch := l.text[l.pos]

	switch ch {
	case '+':
		l.pos++
		return NewToken(TokAdd, start)
	case '-':
		l.pos++
		return NewToken(TokSub, start)
	case '*':
		l.pos++
		return NewToken(TokMul, start)
	case '/':
		l.pos++
		return NewToken(TokDiv, start)
	case '(':
		l.pos++
		return NewToken(TokLParen, start)
	case ')':
		l.pos++
		return NewToken(TokRParen, start)
	default:
		if isDigit(ch) {
			return l.readInt(start)
		}
	}

	l.pos++
	return NewToken(TokInvalid, start)
}

func (l *Lexer) readInt(start uint32) Token {
	for l.pos < uint32(len(l.text)) && isDigit(l.text[l.pos]) {
		l.pos++
	}
	return NewToken(TokInt, start)
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func Lex(text []byte) []Token {
	lexer := NewLexer(text)
	var tokens []Token

	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type() == TokEOF {
			break
		}
	}

	return tokens
}
