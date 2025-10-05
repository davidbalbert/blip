package parse

import (
	"strings"

	"github.com/davidbalbert/blip/lex"
)

// Grammar:
//   Expression → AddExpr
//   AddExpr    → MulExpr (('+' | '-') MulExpr)*
//   MulExpr    → Primary (('*' | '/') Primary)*
//   Primary    → INT | '(' AddExpr ')'

func assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

type NodeType uint8

const (
	NodeInvalid    NodeType = iota
	NodeInt                 // leaf: integer literal
	NodeAdd                 // binary: left + right (childCount = 2)
	NodeSub                 // binary: left - right (childCount = 2)
	NodeMul                 // binary: left * right (childCount = 2)
	NodeDiv                 // binary: left / right (childCount = 2)
	NodeExpression          // root expression node (bracketing)
)

type Node struct {
	Type         NodeType
	TokenID      int // index into tokens slice
	SubtreeStart int // index into nodes slice (first child in postorder)
	HasError     bool
}

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

type ErrorList []error

func (e ErrorList) Error() string {
	messages := make([]string, len(e))
	for i, err := range e {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "\n")
}

func (e ErrorList) Err() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

type Parser struct {
	tokens []lex.Token
	tokIdx int
	nodes  []Node
	errors ErrorList
}

func newParser(tokens []lex.Token) *Parser {
	return &Parser{
		tokens: tokens,
		tokIdx: 0,
		nodes:  make([]Node, 0, len(tokens)),
	}
}

func (p *Parser) currentToken() lex.Token {
	return p.tokens[p.tokIdx]
}

func (p *Parser) currentTokenType() lex.TokenType {
	return p.currentToken().Type()
}

func (p *Parser) consume() lex.Token {
	assert(p.tokIdx < len(p.tokens), "consuming past EOF")
	token := p.currentToken()
	p.tokIdx++
	return token
}

func (p *Parser) addLeafNode(nodeType NodeType, tokenIndex int) {
	p.nodes = append(p.nodes, Node{
		Type:         nodeType,
		TokenID:      tokenIndex,
		SubtreeStart: len(p.nodes), // leaf nodes point to themselves
		HasError:     false,
	})
}

func (p *Parser) addNode(nodeType NodeType, tokenIndex int, subtreeStart int) {
	p.nodes = append(p.nodes, Node{
		Type:         nodeType,
		TokenID:      tokenIndex,
		SubtreeStart: subtreeStart,
		HasError:     false,
	})
}

func (p *Parser) parseAddExpr() {
	subtreeStart := len(p.nodes)

	// Add bracketing node
	exprIdx := p.tokIdx
	p.addLeafNode(NodeExpression, exprIdx)

	// Parse first multiplicative
	p.parseMulExpr()

	// Parse additional multiplicative expressions with operators
	for p.currentTokenType() == lex.TokAdd || p.currentTokenType() == lex.TokSub {
		leftStart := subtreeStart + 1 // Start after bracketing node

		opType := p.currentTokenType()
		opToken := p.tokIdx
		p.consume() // consume operator

		// Parse right operand
		p.parseMulExpr()

		// Emit operator node in postorder (children already emitted)
		if opType == lex.TokAdd {
			p.addNode(NodeAdd, opToken, leftStart)
		} else {
			p.addNode(NodeSub, opToken, leftStart)
		}
	}
}

func (p *Parser) parseMulExpr() {
	leftStart := len(p.nodes)

	// Parse first primary
	p.parseInt()

	// Parse additional primary expressions with operators
	for p.currentTokenType() == lex.TokMul || p.currentTokenType() == lex.TokDiv {
		opType := p.currentTokenType()
		opToken := p.tokIdx
		p.consume() // consume operator

		// Parse right operand
		p.parseInt()

		// Emit operator node in postorder (children already emitted)
		if opType == lex.TokMul {
			p.addNode(NodeMul, opToken, leftStart)
		} else {
			p.addNode(NodeDiv, opToken, leftStart)
		}
	}
}

func (p *Parser) parseInt() {
	if p.currentTokenType() == lex.TokInt {
		tokenIndex := p.tokIdx
		p.consume()
		p.addLeafNode(NodeInt, tokenIndex)
	} else if p.currentTokenType() == lex.TokLParen {
		// Parse parenthesized expression
		p.consume() // consume '('

		// Recursively parse the inner expression (additive level)
		p.parseAddExpr()

		// Consume closing paren (should be there, but handle error if missing)
		if p.currentTokenType() == lex.TokRParen {
			p.consume() // consume ')'
		} else {
			// Error: missing closing paren - emit invalid node
			tokenIndex := p.tokIdx
			p.consume()
			node := Node{
				Type:         NodeInvalid,
				TokenID:      tokenIndex,
				SubtreeStart: len(p.nodes),
				HasError:     true,
			}
			p.nodes = append(p.nodes, node)
		}
	} else {
		// Error case - emit invalid node and consume one token to avoid infinite loop
		tokenIndex := p.tokIdx
		p.consume()
		node := Node{
			Type:         NodeInvalid,
			TokenID:      tokenIndex,
			SubtreeStart: len(p.nodes),
			HasError:     true,
		}
		p.nodes = append(p.nodes, node)
	}
}

func Parse(tokens []lex.Token) ([]Node, error) {
	parser := newParser(tokens)

	if len(tokens) == 0 || (len(tokens) == 1 && tokens[0].Type() == lex.TokEOF) {
		return nil, ErrorList{&ParseError{Message: "expected expression"}}
	}

	parser.parseAddExpr()

	err := parser.errors.Err()
	if err != nil {
		return nil, err
	}
	return parser.nodes, nil
}
