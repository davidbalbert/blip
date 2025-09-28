package arm64

import "encoding/binary"

type Insn uint32

type Generator struct {
	Instructions []Insn
}

func (g *Generator) MovImm(reg int, imm int) {
	// MOV x0, #imm - encoded as MOVZ x0, #imm, lsl #0
	// Format: sf=1, opc=10, hw=00, imm16=imm, Rd=reg
	// 1|1|0|10|0|00|imm16|reg
	insn := Insn(0xD2800000 | (uint32(imm&0xFFFF) << 5) | uint32(reg))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) AddImm(rd, rn int, imm int) {
	// ADD x0, x0, #imm
	// Format: sf=1, op=0, S=0, sh=0, imm12=imm, Rn=rn, Rd=rd
	// 1|0|0|10001|0|0|imm12|rn|rd
	insn := Insn(0x91000000 | (uint32(imm&0xFFF) << 10) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) SubImm(rd, rn int, imm int) {
	// SUB x0, x0, #imm
	// Format: sf=1, op=1, S=0, sh=0, imm12=imm, Rn=rn, Rd=rd
	// 1|1|0|10001|0|0|imm12|rn|rd
	insn := Insn(0xD1000000 | (uint32(imm&0xFFF) << 10) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) MovImm16(reg int, imm int) {
	// MOV x16, #imm for syscall number
	insn := Insn(0xD2800000 | (uint32(imm&0xFFFF) << 5) | uint32(reg))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) SVC(imm int) {
	// SVC #imm - supervisor call
	// Format: 11010100|000|imm16|00000
	insn := Insn(0xD4000001 | (uint32(imm&0xFFFF) << 5))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) Bytes() []byte {
	result := make([]byte, len(g.Instructions)*4)
	for i, insn := range g.Instructions {
		binary.LittleEndian.PutUint32(result[i*4:], uint32(insn))
	}
	return result
}
