package parse

import "github.com/davidbalbert/blip/lex"

type NodeKind uint8

const (
	NodeInvalid    NodeKind = iota
	NodeInt                 // leaf: integer literal
	NodeAdd                 // binary: left + right (childCount = 2)
	NodeSub                 // binary: left - right (childCount = 2)
	NodeMul                 // binary: left * right (childCount = 2)
	NodeDiv                 // binary: left / right (childCount = 2)
	NodeExpression          // root expression node (bracketing)
)

type Node struct {
	Kind         NodeKind
	Token        int // index into tokens slice
	SubtreeStart int // index into nodes slice (first child in postorder)
	HasError     bool
}

type Parser struct {
	tokens []lex.Token
	pos    int    // current token index
	nodes  []Node // postorder storage
}

func newParser(tokens []lex.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		nodes:  make([]Node, 0, len(tokens)), // pre-size for efficiency
	}
}

func (p *Parser) currentToken() lex.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1] // EOF token
	}
	return p.tokens[p.pos]
}

func (p *Parser) currentTokenType() lex.TokenType {
	return p.currentToken().Type()
}

func (p *Parser) consume() lex.Token {
	token := p.currentToken()
	if p.pos < len(p.tokens)-1 { // don't advance past EOF
		p.pos++
	}
	return token
}

func (p *Parser) addLeafNode(kind NodeKind, tokenIndex int) {
	p.nodes = append(p.nodes, Node{
		Kind:         kind,
		Token:        tokenIndex,
		SubtreeStart: len(p.nodes), // points to itself for leaves
		HasError:     false,
	})
}

func (p *Parser) addNode(kind NodeKind, tokenIndex int, subtreeStart int) {
	p.nodes = append(p.nodes, Node{
		Kind:         kind,
		Token:        tokenIndex,
		SubtreeStart: subtreeStart,
		HasError:     false,
	})
}

// Parser for arithmetic expressions with proper precedence
// Grammar:
//
//	addExpr = mulExpr (('+' | '-') mulExpr)*
//	mulExpr = intExpr (('*' | '/') intExpr)*
//	int = INT
func (p *Parser) parseAddExpr() {
	subtreeStart := len(p.nodes)

	// Add bracketing node
	expressionToken := p.pos
	p.addLeafNode(NodeExpression, expressionToken)

	// Parse first multiplicative
	p.parseMulExpr()

	// Parse additional multiplicative expressions with operators
	for p.currentTokenType() == lex.TokAdd || p.currentTokenType() == lex.TokSub {
		opType := p.currentTokenType()
		opToken := p.pos
		p.consume() // consume operator

		// Parse right operand
		p.parseMulExpr()

		// Emit operator node in postorder (children already emitted)
		if opType == lex.TokAdd {
			p.addNode(NodeAdd, opToken, subtreeStart)
		} else {
			p.addNode(NodeSub, opToken, subtreeStart)
		}

		// Update subtreeStart for next operator
		subtreeStart = len(p.nodes) - 1
	}
}

func (p *Parser) parseMulExpr() {
	subtreeStart := len(p.nodes)

	// Parse first primary
	p.parseInt()

	// Parse additional primary expressions with operators
	for p.currentTokenType() == lex.TokMul || p.currentTokenType() == lex.TokDiv {
		opType := p.currentTokenType()
		opToken := p.pos
		p.consume() // consume operator

		// Parse right operand
		p.parseInt()

		// Emit operator node in postorder (children already emitted)
		if opType == lex.TokMul {
			p.addNode(NodeMul, opToken, subtreeStart)
		} else {
			p.addNode(NodeDiv, opToken, subtreeStart)
		}

		// Update subtreeStart for next operator
		subtreeStart = len(p.nodes) - 1
	}
}

func (p *Parser) parseInt() {
	if p.currentTokenType() == lex.TokInt {
		tokenIndex := p.pos
		p.consume()
		p.addLeafNode(NodeInt, tokenIndex)
	} else {
		// Error case - emit invalid node and consume one token to avoid infinite loop
		tokenIndex := p.pos
		p.consume()
		node := Node{
			Kind:         NodeInvalid,
			Token:        tokenIndex,
			SubtreeStart: len(p.nodes),
			HasError:     true,
		}
		p.nodes = append(p.nodes, node)
	}
}

// Parse parses the tokens into a parse tree
func Parse(tokens []lex.Token) []Node {
	parser := newParser(tokens)

	if len(tokens) == 0 || (len(tokens) == 1 && tokens[0].Type() == lex.TokEOF) {
		return parser.nodes
	}

	parser.parseAddExpr()

	return parser.nodes
}
