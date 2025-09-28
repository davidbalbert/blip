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

func (g *Generator) Add(rd, rn, rm int) {
	// ADD xd, xn, xm
	// Format: sf=1, op=0, S=0, shift=00, Rm=rm, imm6=000000, Rn=rn, Rd=rd
	// 1|0|0|01011|00|0|rm|000000|rn|rd
	insn := Insn(0x8B000000 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) SubImm(rd, rn int, imm int) {
	// SUB x0, x0, #imm
	// Format: sf=1, op=1, S=0, sh=0, imm12=imm, Rn=rn, Rd=rd
	// 1|1|0|10001|0|0|imm12|rn|rd
	insn := Insn(0xD1000000 | (uint32(imm&0xFFF) << 10) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) Sub(rd, rn, rm int) {
	// SUB xd, xn, xm
	// Format: sf=1, op=1, S=0, shift=00, Rm=rm, imm6=000000, Rn=rn, Rd=rd
	// 1|1|0|01011|00|0|rm|000000|rn|rd
	insn := Insn(0xCB000000 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) MovImm16(reg int, imm int) {
	// MOV x16, #imm for syscall number
	insn := Insn(0xD2800000 | (uint32(imm&0xFFFF) << 5) | uint32(reg))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) Mul(rd, rn, rm int) {
	// MUL xd, xn, xm (multiply)
	// Format: sf=1, op54=00, op31=11010110, Rm=rm, op15=011111, Rn=rn, Rd=rd
	// 1|00|11010110|rm|011111|rn|rd
	insn := Insn(0x9B007C00 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) Div(rd, rn, rm int) {
	// SDIV xd, xn, xm (signed divide)
	// Format: sf=1, op1=0, S=0, op2=11010110, Rm=rm, op3=000011, Rn=rn, Rd=rd
	// 1|0|0|11010110|rm|000011|rn|rd
	insn := Insn(0x9AC00C00 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) SVC(imm int) {
	// SVC #imm - supervisor call
	// Format: 11010100|000|imm16|00000
	insn := Insn(0xD4000001 | (uint32(imm&0xFFFF) << 5))
	g.Instructions = append(g.Instructions, insn)
}

func (g *Generator) PushReg(reg int) {
	// sub sp, sp, #16; str xN, [sp]
	// SUB sp, sp, #16
	insn1 := Insn(0xD10043FF) // sub sp, sp, #16
	g.Instructions = append(g.Instructions, insn1)
	// STR xN, [sp]
	insn2 := Insn(0xF90003E0 | uint32(reg)) // str xN, [sp]
	g.Instructions = append(g.Instructions, insn2)
}

func (g *Generator) PopReg(reg int) {
	// ldr xN, [sp]; add sp, sp, #16
	// LDR xN, [sp]
	insn1 := Insn(0xF94003E0 | uint32(reg)) // ldr xN, [sp]
	g.Instructions = append(g.Instructions, insn1)
	// ADD sp, sp, #16
	insn2 := Insn(0x910043FF) // add sp, sp, #16
	g.Instructions = append(g.Instructions, insn2)
}

func (g *Generator) Bytes() []byte {
	result := make([]byte, len(g.Instructions)*4)
	for i, insn := range g.Instructions {
		binary.LittleEndian.PutUint32(result[i*4:], uint32(insn))
	}
	return result
}
