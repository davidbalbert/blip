package obj

import (
	"strconv"

	"github.com/davidbalbert/blip/lex"
	"github.com/davidbalbert/blip/obj/arm64"
	"github.com/davidbalbert/blip/parse"
)

// Codegen generates ARM64 machine code from a parse tree
func Codegen(nodes []parse.Node, tokens []lex.Token, text []byte) []byte {
	gen := arm64.Generator{}

	ctx := &codegenContext{
		gen:    &gen,
		nodes:  nodes,
		tokens: tokens,
		text:   text,
		depth:  0,
	}

	ctx.walkNodes()

	// Add exit syscall
	gen.MovImm16(16, 1) // mov x16, #1 (sys_exit)
	gen.SVC(0x80)       // svc #0x80

	return gen.Bytes()
}

type codegenContext struct {
	gen    *arm64.Generator
	nodes  []parse.Node
	tokens []lex.Token
	text   []byte
	depth  int // operand stack depth
}

func (ctx *codegenContext) walkNodes() {
	for _, node := range ctx.nodes {
		ctx.processNode(node)
	}
}

func (ctx *codegenContext) processNode(node parse.Node) {
	switch node.Kind {
	case parse.NodeExpression:
		// Bracketing node - no code generated
		return

	case parse.NodeInt:
		token := ctx.tokens[node.Token]
		value := ctx.getIntValue(token)

		if ctx.depth == 0 {
			// First operand goes to x0
			ctx.gen.MovImm(0, value)
		} else {
			// Push current x0 to stack, load new value
			ctx.gen.PushReg(0)
			ctx.gen.MovImm(0, value)
		}
		ctx.depth++

	case parse.NodeAdd:
		// Pop left operand into x1, right is in x0
		ctx.gen.PopReg(1)    // left operand
		ctx.gen.Add(0, 1, 0) // x0 = x1 + x0
		ctx.depth--

	case parse.NodeSub:
		// Pop left operand into x1, right is in x0
		ctx.gen.PopReg(1)    // left operand
		ctx.gen.Sub(0, 1, 0) // x0 = x1 - x0
		ctx.depth--

	case parse.NodeMul:
		// Pop left operand into x1, right is in x0
		ctx.gen.PopReg(1)    // left operand
		ctx.gen.Mul(0, 1, 0) // x0 = x1 * x0
		ctx.depth--

	case parse.NodeDiv:
		// Pop left operand into x1, right is in x0
		ctx.gen.PopReg(1)    // left operand
		ctx.gen.Div(0, 1, 0) // x0 = x1 / x0
		ctx.depth--
	}
}

func (ctx *codegenContext) getIntValue(token lex.Token) int {
	start := token.Pos()
	end := start
	for end < uint32(len(ctx.text)) && isDigit(ctx.text[end]) {
		end++
	}
	value, _ := strconv.Atoi(string(ctx.text[start:end]))
	return value
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
