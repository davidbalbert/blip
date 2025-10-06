package parse

import (
	"strings"

	"github.com/davidbalbert/blip/lex"
)

// Grammar:
//   Expr → AddExpr
//   AddExpr    → MulExpr (('+' | '-') MulExpr)*
//   MulExpr    → Primary (('*' | '/') Primary)*
//   Primary    → INT | '(' AddExpr ')'

type NodeType uint8

const (
	NodeInvalid        NodeType = iota
	NodeInt                     // leaf: integer literal
	NodeAdd                     // binary: left + right (childCount = 2)
	NodeSub                     // binary: left - right (childCount = 2)
	NodeMul                     // binary: left * right (childCount = 2)
	NodeDiv                     // binary: left / right (childCount = 2)
	NodeParenExprStart          // bracketing: '(' introducer (first child)
	NodeParenExpr               // parenthesized expression root
	NodeExpr                    // root expression node
)

func (t NodeType) String() string {
	switch t {
	case NodeInvalid:
		return "Invalid"
	case NodeInt:
		return "Int"
	case NodeAdd:
		return "Add"
	case NodeSub:
		return "Sub"
	case NodeMul:
		return "Mul"
	case NodeDiv:
		return "Div"
	case NodeParenExprStart:
		return "ParenExprStart"
	case NodeParenExpr:
		return "ParenExpr"
	case NodeExpr:
		return "Expr"
	default:
		return "Unknown"
	}
}

type Node struct {
	Type         NodeType
	TokenID      int // index into tokens slice
	SubtreeStart int // index into nodes slice (first child in postorder)
	HasError     bool
}

type ParseError struct {
	message     string
	sourceLine  string
	errorOffset uint32
}

func (e *ParseError) Error() string {
	caret := strings.Repeat(" ", int(e.errorOffset)) + "^"
	return e.sourceLine + "\n" + caret + "\n" + e.message
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
	tokens     []lex.Token
	lineStarts []uint32
	lineIDs    []int
	source     []byte
	tokenIndex int
	nodes      []Node
	errors     ErrorList
}

func newParser(tokens lex.Tokens, source []byte) *Parser {
	return &Parser{
		tokens:     tokens.Tokens,
		lineStarts: tokens.LineStarts,
		lineIDs:    tokens.LineIDs,
		source:     source,
		tokenIndex: 0,
		nodes:      make([]Node, 0, len(tokens.Tokens)),
	}
}

func (p *Parser) currentToken() lex.Token {
	return p.tokens[p.tokenIndex]
}

func (p *Parser) currentTokenType() lex.TokenType {
	return p.currentToken().Type()
}

func (p *Parser) consume() lex.Token {
	token := p.currentToken()
	p.tokenIndex++
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

func (p *Parser) addError(message string, tokenIndex int) {
	pos := p.tokens[tokenIndex].Pos()
	lineID := p.lineIDs[tokenIndex]
	lineStart := p.lineStarts[lineID]

	var lineEnd uint32
	if lineID+1 < len(p.lineStarts) {
		lineEnd = p.lineStarts[lineID+1] - 1
	} else {
		lineEnd = uint32(len(p.source))
	}

	sourceLine := string(p.source[lineStart:lineEnd])

	p.errors = append(p.errors, &ParseError{
		message:     message,
		sourceLine:  sourceLine,
		errorOffset: pos - lineStart,
	})
}

func (p *Parser) parseAddExpr() {
	// Parse first multiplicative
	p.parseMulExpr()

	// Parse additional multiplicative expressions with operators
	for p.currentTokenType() == lex.TokAdd || p.currentTokenType() == lex.TokSub {
		leftStart := len(p.nodes) - 1

		opType := p.currentTokenType()
		opToken := p.tokenIndex
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
	p.parsePrimary()

	// Parse additional primary expressions with operators
	for p.currentTokenType() == lex.TokMul || p.currentTokenType() == lex.TokDiv {
		opType := p.currentTokenType()
		opToken := p.tokenIndex
		p.consume() // consume operator

		// Parse right operand
		p.parsePrimary()

		// Emit operator node in postorder (children already emitted)
		if opType == lex.TokMul {
			p.addNode(NodeMul, opToken, leftStart)
		} else {
			p.addNode(NodeDiv, opToken, leftStart)
		}
	}
}

func (p *Parser) parsePrimary() {
	if p.currentTokenType() == lex.TokInt {
		tokenIndex := p.tokenIndex
		p.consume()
		p.addLeafNode(NodeInt, tokenIndex)
	} else if p.currentTokenType() == lex.TokLParen {
		subtreeStart := len(p.nodes)

		lparenToken := p.tokenIndex
		p.consume() // consume '('

		p.addLeafNode(NodeParenExprStart, lparenToken)

		p.parseAddExpr()

		if p.currentTokenType() == lex.TokRParen {
			rparenToken := p.tokenIndex
			p.consume() // consume ')'
			p.addNode(NodeParenExpr, rparenToken, subtreeStart)
		} else {
			tokenIndex := p.tokenIndex
			p.addError("expected ')' or operator", tokenIndex)
			if p.currentTokenType() != lex.TokEOF {
				p.consume()
			}
			node := Node{
				Type:         NodeInvalid,
				TokenID:      tokenIndex,
				SubtreeStart: len(p.nodes),
				HasError:     true,
			}
			p.nodes = append(p.nodes, node)
		}
	} else {
		tokenIndex := p.tokenIndex
		p.addError("expected expression", tokenIndex)
		if p.currentTokenType() != lex.TokEOF {
			p.consume()
		}
		node := Node{
			Type:         NodeInvalid,
			TokenID:      tokenIndex,
			SubtreeStart: len(p.nodes),
			HasError:     true,
		}
		p.nodes = append(p.nodes, node)
	}
}

func (p *Parser) parseExpr() {
	if p.currentTokenType() == lex.TokEOF {
		p.addError("expected expression", p.tokenIndex)
		return
	}

	subtreeStart := len(p.nodes)

	p.parseAddExpr()

	if p.currentTokenType() != lex.TokEOF {
		p.addError("expected operator or end of expression", p.tokenIndex)
		if p.currentTokenType() != lex.TokEOF {
			p.consume()
		}
		node := Node{
			Type:         NodeInvalid,
			TokenID:      p.tokenIndex,
			SubtreeStart: len(p.nodes),
			HasError:     true,
		}
		p.nodes = append(p.nodes, node)
		return
	}

	eofToken := p.tokenIndex
	p.consume() // consume EOF

	p.addNode(NodeExpr, eofToken, subtreeStart)
}

func Parse(tokens lex.Tokens, source []byte) ([]Node, error) {
	parser := newParser(tokens, source)

	parser.parseExpr()
	return parser.nodes, parser.errors.Err()
}
