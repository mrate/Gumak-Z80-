package z80

import "mutex/gumak/log"

// TEST
func testBit0(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 0) }
func testBit1(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 1) }
func testBit2(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 2) }
func testBit3(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 3) }
func testBit4(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 4) }
func testBit5(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 5) }
func testBit6(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 6) }
func testBit7(cpu *CPU, ans uint8) { Alu_BIT(cpu, ans, 7) }

// SET
func setBit0(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 0) }
func setBit1(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 1) }
func setBit2(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 2) }
func setBit3(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 3) }
func setBit4(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 4) }
func setBit5(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 5) }
func setBit6(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 6) }
func setBit7(cpu *CPU, ans uint8) uint8 { return Alu_SET(cpu, ans, 7) }

// RESET
func resetBit0(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 0) }
func resetBit1(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 1) }
func resetBit2(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 2) }
func resetBit3(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 3) }
func resetBit4(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 4) }
func resetBit5(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 5) }
func resetBit6(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 6) }
func resetBit7(cpu *CPU, ans uint8) uint8 { return Alu_RES(cpu, ans, 7) }

func Instr_0xCB(cpu *CPU) int {
	op2 := FetchInstruction(cpu)
	r := op2 & 0b00000111

	test := op2 >= 64 && op2 < 128
	var wOp func(cpu *CPU, val uint8) uint8
	var tOp func(cpu *CPU, val uint8)

	switch {
	case op2 < 8:
		log.Trace(2, "RLC")
		wOp = Alu_RLC
	case op2 < 16:
		log.Trace(2, "RRC")
		wOp = Alu_RRC
	case op2 < 24:
		log.Trace(2, "RL")
		wOp = Alu_RL
	case op2 < 32:
		log.Trace(2, "RR")
		wOp = Alu_RR
	case op2 < 40:
		log.Trace(2, "SLA")
		wOp = Alu_SLA
	case op2 < 48:
		log.Trace(2, "SRA")
		wOp = Alu_SRA
	case op2 < 56:
		log.Trace(2, "SLS")
		wOp = Alu_SLS
	case op2 < 64:
		log.Trace(2, "SRL")
		wOp = Alu_SRL
	case op2 < 72:
		// TEST BIT
		log.Trace(2, "BIT 0, ")
		tOp = testBit0
	case op2 < 80:
		log.Trace(2, "BIT 1, ")
		tOp = testBit1
	case op2 < 88:
		log.Trace(2, "BIT 2, ")
		tOp = testBit2
	case op2 < 96:
		log.Trace(2, "BIT 3, ")
		tOp = testBit3
	case op2 < 104:
		log.Trace(2, "BIT 4, ")
		tOp = testBit4
	case op2 < 112:
		log.Trace(2, "BIT 5, ")
		tOp = testBit5
	case op2 < 120:
		log.Trace(2, "BIT 6, ")
		tOp = testBit6
	case op2 < 128:
		log.Trace(2, "BIT 7, ")
		tOp = testBit7
	case op2 < 136:
		// RES
		log.Trace(2, "RES 0, ")
		wOp = resetBit0
	case op2 < 144:
		log.Trace(2, "RES 1, ")
		wOp = resetBit1
	case op2 < 152:
		log.Trace(2, "RES 2, ")
		wOp = resetBit2
	case op2 < 160:
		log.Trace(2, "RES 3, ")
		wOp = resetBit3
	case op2 < 168:
		log.Trace(2, "RES 4, ")
		wOp = resetBit4
	case op2 < 176:
		log.Trace(2, "RES 5, ")
		wOp = resetBit5
	case op2 < 184:
		log.Trace(2, "RES 6, ")
		wOp = resetBit6
	case op2 < 192:
		log.Trace(2, "RES 7, ")
		wOp = resetBit7
	case op2 < 200:
		log.Trace(2, "SET 0, ")
		wOp = setBit0
	case op2 < 208:
		log.Trace(2, "SET 1, ")
		wOp = setBit1
	case op2 < 216:
		log.Trace(2, "SET 2, ")
		wOp = setBit2
	case op2 < 224:
		log.Trace(2, "SET 3, ")
		wOp = setBit3
	case op2 < 232:
		log.Trace(2, "SET 4, ")
		wOp = setBit4
	case op2 < 240:
		log.Trace(2, "SET 5, ")
		wOp = setBit5
	case op2 < 248:
		log.Trace(2, "SET 6, ")
		wOp = setBit6
	default:
		log.Trace(2, "SET 7, ")
		wOp = setBit7
	}

	if r == 0b110 {
		value := MEM_HL(cpu)
		log.Trace(2, "(HL)")

		if test {
			tOp(cpu, value)
		} else {
			MEM_HL_W(cpu, wOp(cpu, value))
		}

		return 12
	} else {
		var reg *uint8
		switch r {
		case 0b000:
			reg = &cpu.Reg.B
		case 0b001:
			reg = &cpu.Reg.C
		case 0b010:
			reg = &cpu.Reg.D
		case 0b011:
			reg = &cpu.Reg.E
		case 0b100:
			reg = &cpu.Reg.H
		case 0b101:
			reg = &cpu.Reg.L
		case 0b111:
			reg = &cpu.Reg.A
		}

		log.Trace(2, "%s", cpu.Reg.Name(reg))

		if test {
			tOp(cpu, *reg)
		} else {
			*reg = wOp(cpu, *reg)
		}

		return 8
	}
}

func BIT_IX_IY_cb(cpu *CPU, dst *uint16, d int, bit uint8) int {
	addr := uint16(int(*dst) + d)

	log.Trace(2, "BIT %d, (%s+%d)", bit, cpu.Reg.Name16(dst), d)
	Alu_BIT(cpu, MemoryRead(cpu, addr), bit)
	return 20
}

func SET_IX_IY_cb(cpu *CPU, dst *uint16, d int, bit uint8) int {
	addr := uint16(int(*dst) + d)

	log.Trace(2, "SET %d, (%s+%d)", bit, cpu.Reg.Name16(dst), d)
	MemoryWrite(cpu, addr, Alu_SET(cpu, MemoryRead(cpu, addr), bit))
	return 23
}

func RES_IX_IY_cb(cpu *CPU, dst *uint16, d int, bit uint8) int {
	addr := uint16(int(*dst) + d)

	log.Trace(2, "RES %d, (%s+%d)", bit, cpu.Reg.Name16(dst), d)
	MemoryWrite(cpu, addr, Alu_RES(cpu, MemoryRead(cpu, addr), bit))
	return 23
}

func RLD(cpu *CPU) int {
	value := MEM_HL(cpu)
	newVal := uint8((value << 4) | cpu.Reg.A&0x0F)
	MEM_HL_W(cpu, newVal)

	cpu.Reg.A = (cpu.Reg.A & 0xF0) | (value >> 4)

	cpu.SetFlag(FLAG_Z, cpu.Reg.A == 0)
	cpu.SetFlag(FLAG_S, (cpu.Reg.A&0x80) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, Alu_ParityEven(cpu.Reg.A))
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)

	log.Trace(2, "RLD")
	return 18
}

func RRD(cpu *CPU) int {
	value := MEM_HL(cpu)
	newVal := uint8(((cpu.Reg.A & 0x0F) << 4) | (value >> 4))
	MEM_HL_W(cpu, newVal)

	cpu.Reg.A = (cpu.Reg.A & 0xF0) | (value & 0x0F)

	cpu.SetFlag(FLAG_Z, cpu.Reg.A == 0)
	cpu.SetFlag(FLAG_S, (cpu.Reg.A&0x80) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, Alu_ParityEven(cpu.Reg.A))
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)

	log.Trace(2, "RRD")
	return 18
}
