package main

import (
	"strconv"

	"github.com/davidbalbert/blip/lex"
	"github.com/davidbalbert/blip/obj/arm64"
	"github.com/davidbalbert/blip/parse"
)

// codegen generates ARM64 machine code from a parse tree
func codegen(nodes []parse.ParseNode, tokens []lex.Token, text []byte) []byte {
	gen := arm64.Generator{}
	
	// Walk the parse tree and generate code
	// For now, we'll traverse the nodes and handle the arithmetic expression
	ctx := &codegenContext{
		gen:    &gen,
		nodes:  nodes,
		tokens: tokens,
		text:   text,
		firstValue: true,
	}
	
	ctx.walkNodes()
	
	// Add exit syscall
	gen.MovImm16(16, 1) // mov x16, #1 (sys_exit)
	gen.SVC(0x80)       // svc #0x80
	
	return gen.Bytes()
}

type codegenContext struct {
	gen        *arm64.Generator
	nodes      []parse.ParseNode
	tokens     []lex.Token
	text       []byte
	firstValue bool
	nodeIndex  int
}

func (ctx *codegenContext) walkNodes() {
	for ctx.nodeIndex < len(ctx.nodes) {
		ctx.processNode()
	}
}

func (ctx *codegenContext) processNode() {
	if ctx.nodeIndex >= len(ctx.nodes) {
		return
	}
	
	node := ctx.nodes[ctx.nodeIndex]
	ctx.nodeIndex++
	
	switch node.Kind {
	case parse.NodeExpression:
		// Bracketing node - just process children that follow
		return
		
	case parse.NodeInt:
		// Generate code for integer literal
		token := ctx.tokens[node.Token]
		value := ctx.getIntValue(token)
		
		if ctx.firstValue {
			ctx.gen.MovImm(0, value) // mov x0, #value
			ctx.firstValue = false
		}
		
	case parse.NodeAdd:
		// Binary add - right operand is on stack, left operand in x0
		// Process right operand (it's the previous node due to postorder)
		rightNode := ctx.nodes[ctx.nodeIndex-2] // -1 for current, -1 for right operand
		if rightNode.Kind == parse.NodeInt {
			rightToken := ctx.tokens[rightNode.Token]
			rightValue := ctx.getIntValue(rightToken)
			ctx.gen.AddImm(0, 0, rightValue) // add x0, x0, #rightValue
		}
		
	case parse.NodeSub:
		// Binary sub - similar to add
		rightNode := ctx.nodes[ctx.nodeIndex-2]
		if rightNode.Kind == parse.NodeInt {
			rightToken := ctx.tokens[rightNode.Token]
			rightValue := ctx.getIntValue(rightToken)
			ctx.gen.SubImm(0, 0, rightValue) // sub x0, x0, #rightValue
		}
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
