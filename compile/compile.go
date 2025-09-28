package compile

import (
	"github.com/davidbalbert/blip/lex"
	"github.com/davidbalbert/blip/macho"
	"github.com/davidbalbert/blip/obj"
	"github.com/davidbalbert/blip/parse"
)

// Compile compiles source code to Mach-O object file bytes
func Compile(content []byte) []byte {
	tokens := lex.Lex(content)
	nodes := parse.Parse(tokens)
	text := obj.Codegen(nodes, tokens, content)

	return macho.Generate(text)
}
