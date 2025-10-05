package parse

import (
	"fmt"
	"testing"

	"github.com/davidbalbert/blip/lex"
)

func tokenText(tokens []lex.Token, tokenID int, source string) string {
	token := tokens[tokenID]
	start := token.Pos()

	var end uint32
	if tokenID+1 < len(tokens) {
		end = tokens[tokenID+1].Pos()
	} else {
		end = uint32(len(source))
	}

	text := source[start:end]
	for len(text) > 0 && (text[len(text)-1] == ' ' || text[len(text)-1] == '\t' || text[len(text)-1] == '\n' || text[len(text)-1] == '\r') {
		text = text[:len(text)-1]
	}

	return text
}

func sexpr(nodes []Node, tokens []lex.Token, source string) string {
	if len(nodes) == 0 {
		return "()"
	}

	stack := []string{}

	for i, node := range nodes {
		tokenText := tokenText(tokens, node.TokenID, source)

		switch node.Type {
		case NodeInt:
			stack = append(stack, tokenText)
		case NodeAdd, NodeSub, NodeMul, NodeDiv:
			if len(stack) < 2 {
				return "<malformed tree>"
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			op := ""
			switch node.Type {
			case NodeAdd:
				op = "add"
			case NodeSub:
				op = "sub"
			case NodeMul:
				op = "mul"
			case NodeDiv:
				op = "div"
			}
			stack = append(stack, fmt.Sprintf("(%s %s %s)", op, left, right))
		case NodeParenExprStart:
			stack = append(stack, "paren-start")
		case NodeParenExpr:
			if len(stack) < 2 {
				return "<malformed tree>"
			}
			innerExpr := stack[len(stack)-1]
			lparen := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, fmt.Sprintf("(paren %s %s)", lparen, innerExpr))
		case NodeExpr:
			if len(stack) < 1 {
				return "<malformed tree>"
			}
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, fmt.Sprintf("(expr %s)", child))
		case NodeInvalid:
			stack = append(stack, fmt.Sprintf("(invalid %s)", tokenText))
		default:
			stack = append(stack, fmt.Sprintf("(unknown-%d %s)", node.Type, tokenText))
		}

		// Debug: show stack state
		_ = i
	}

	if len(stack) != 1 {
		return fmt.Sprintf("<stack size %d, expected 1: %v>", len(stack), stack)
	}

	return stack[0]
}

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"simple_int", "42", "(expr 42)"},
		{"add", "1 + 2", "(expr (add 1 2))"},
		{"sub", "5 - 3", "(expr (sub 5 3))"},
		{"mul", "3 * 4", "(expr (mul 3 4))"},
		{"div", "8 / 2", "(expr (div 8 2))"},
		{"precedence", "2 + 3 * 4", "(expr (add 2 (mul 3 4)))"},
		{"left_assoc_add", "1 + 2 + 3", "(expr (add (add 1 2) 3))"},
		{"left_assoc_mul", "2 * 3 * 4", "(expr (mul (mul 2 3) 4))"},
		{"complex", "10 + 2 * 3 - 8 / 4", "(expr (sub (add 10 (mul 2 3)) (div 8 4)))"},
		{"parens_basic", "(2 + 3) * 4", "(expr (mul (paren paren-start (add 2 3)) 4))"},
		{"parens_nested", "((2 + 3) * 4) - 5", "(expr (sub (paren paren-start (mul (paren paren-start (add 2 3)) 4)) 5))"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lex.Lex([]byte(tc.source))
			nodes, err := Parse(tokens)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}

			// in a valid parse tree there's a 1:1 mapping from nodes to tokens
			if len(nodes) != len(tokens) {
				t.Errorf("should be equal: len(nodes)=%d, len(tokens)=%d", len(nodes), len(tokens))
			}

			seen := make(map[int]bool)
			for _, node := range nodes {
				if seen[node.TokenID] {
					t.Errorf("token %d appears in multiple nodes", node.TokenID)
				}
				seen[node.TokenID] = true
			}

			actual := sexpr(nodes, tokens, tc.source)
			if actual != tc.expected {
				t.Errorf("got \"%s\", want \"%s\"", actual, tc.expected)
			}
		})
	}
}
