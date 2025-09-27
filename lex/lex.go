package lex

type TokenType uint8

const (
	TokInvalid TokenType = iota
	TokEOF
	TokInt
	TokAdd
	TokSub
)

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
