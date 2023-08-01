package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

func ADD_A(cpu *CPU, reg *uint8) int {
	Alu_ADD_A(cpu, *reg)

	log.Trace(2, "ADD A, %s", cpu.Reg.Name(reg))
	return 4
}

func ADC_A(cpu *CPU, reg *uint8) int {
	Alu_ADC_A(cpu, *reg)

	log.Trace(2, "ADD C, %s", cpu.Reg.Name(reg))
	return 4
}

func SUB_A(cpu *CPU, reg *uint8) int {
	Alu_SUB_A(cpu, *reg)

	log.Trace(2, "SUB A, %s", cpu.Reg.Name(reg))
	return 4
}

func SUB_A_HL(cpu *CPU) int {
	Alu_SUB_A(cpu, MEM_HL(cpu))

	log.Trace(2, "SUB A, (HL)")
	return 7
}

func INC(cpu *CPU, reg *uint8) int {
	*reg = Alu_INC8(cpu, *reg)

	log.Trace(2, "INC %s", cpu.Reg.Name(reg))
	return 4
}

func UN_INC_HL(cpu *CPU, reg *uint16, low bool) int {
	l, h := helpers.To8(*reg)

	if low {
		l = Alu_INC8(cpu, l)
	} else {
		h = Alu_INC8(cpu, h)
	}

	*reg = helpers.To16(h, l)

	log.Trace(2, "INC %s", cpu.Reg.Name16HL(reg, low))
	return 4
}

func UN_DEC_HL(cpu *CPU, reg *uint16, low bool) int {
	l, h := helpers.To8(*reg)

	if low {
		l = Alu_DEC8(cpu, l)
	} else {
		h = Alu_DEC8(cpu, h)
	}

	*reg = helpers.To16(h, l)

	log.Trace(2, "INC %s", cpu.Reg.Name16HL(reg, low))
	return 4
}

func INC_HL(cpu *CPU) int {
	MEM_HL_W(cpu, Alu_INC8(cpu, MEM_HL(cpu)))

	log.Trace(2, "INC (HL)")
	return 11
}

func DEC(cpu *CPU, reg *uint8) int {
	*reg = Alu_DEC8(cpu, *reg)

	log.Trace(2, "DEC %s", cpu.Reg.Name(reg))
	return 4
}

func DEC_HL(cpu *CPU) int {
	MEM_HL_W(cpu, Alu_DEC8(cpu, MEM_HL(cpu)))

	log.Trace(2, "DEC (HL)")
	return 11
}

func INC16(cpu *CPU, h *uint8, l *uint8) int {
	*l, *h = helpers.To8(helpers.To16(*l, *h) + 1)

	log.Trace(2, "INC %s%s", cpu.Reg.Name(h), cpu.Reg.Name(l))
	return 4
}

func DEC16(cpu *CPU, h *uint8, l *uint8) int {
	*l, *h = helpers.To8(helpers.To16(*l, *h) - 1)

	log.Trace(2, "DEC %s%s", cpu.Reg.Name(h), cpu.Reg.Name(l))
	return 4
}

func DEC_r(cpu *CPU, reg *uint8) int {
	*reg = Alu_DEC8(cpu, *reg)

	log.Trace(2, "DEC %s", cpu.Reg.Name(reg))
	return 4
}

func ADD16(cpu *CPU, rh *uint8, rl *uint8, sh *uint8, sl *uint8) int {
	value := Alu_ADD16(cpu, helpers.To16(*rl, *rh), helpers.To16(*sl, *sh))
	*rl, *rh = helpers.To8(value)

	log.Trace(2, "ADD %s%s, %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl), cpu.Reg.Name(sh), cpu.Reg.Name(sh))
	return 11
}

func ADD_A_n(cpu *CPU, n uint8) int {
	Alu_ADD_A(cpu, n)

	log.Trace(2, "ADD A, $%02x", n)
	return 7
}

func SBC_A(cpu *CPU, reg *uint8) int {
	Alu_SBC_A(cpu, *reg)

	log.Trace(2, "SBC A, %s", cpu.Reg.Name(reg))
	return 4
}

func SBC_A_HL(cpu *CPU) int {
	Alu_SBC_A(cpu, MEM_HL(cpu))

	log.Trace(2, "SBC A, (HL)")
	return 7
}

func AND_A(cpu *CPU, reg *uint8) int {
	Alu_AND_A(cpu, *reg)

	log.Trace(2, "AND A, %s", cpu.Reg.Name(reg))
	return 4
}

func OR_A(cpu *CPU, reg *uint8) int {
	Alu_OR_A(cpu, *reg)

	log.Trace(2, "OR A, %s", cpu.Reg.Name(reg))
	return 4
}

func XOR_A(cpu *CPU, reg *uint8) int {
	Alu_XOR_A(cpu, *reg)

	log.Trace(2, "XOR A, %s", cpu.Reg.Name(reg))
	return 4
}

func CP_A(cpu *CPU, reg *uint8) int {
	Alu_CP_A(cpu, *reg)

	log.Trace(2, "CP %s", cpu.Reg.Name(reg))
	return 4
}

// 16-bin

func ADD_HL_16(cpu *CPU, rh *uint8, rl *uint8) int {
	cpu.Reg.HL_write(Alu_ADD16(cpu, cpu.Reg.HL(), helpers.To16(*rl, *rh)))
	log.Trace(2, "ADD HL, %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 11
}

func ADC_HL_16(cpu *CPU, rh *uint8, rl *uint8) int {
	cpu.Reg.HL_write(Alu_ADC16(cpu, cpu.Reg.HL(), helpers.To16(*rl, *rh)))

	log.Trace(2, "ADC HL, %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 15
}

func SBC_HL_16(cpu *CPU, rh *uint8, rl *uint8) int {
	cpu.Reg.HL_write(Alu_SBC16(cpu, cpu.Reg.HL(), helpers.To16(*rl, *rh)))

	log.Trace(2, "SBC HL, %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 15
}

// IX / IY

func INC_IXIYd(cpu *CPU, reg *uint16) int {
	d := FetchOperand8Compl(cpu)

	addr := uint16(int(*reg) + d)
	MemoryWrite(cpu, addr, Alu_INC8(cpu, MemoryRead(cpu, addr)))

	log.Trace(2, "INC (%s+%d)", cpu.Reg.Name16(reg), d)
	return 23
}

func DEC_IXIYd(cpu *CPU, reg *uint16) int {
	d := FetchOperand8Compl(cpu)

	addr := uint16(int(*reg) + d)
	MemoryWrite(cpu, addr, Alu_DEC8(cpu, MemoryRead(cpu, addr)))

	log.Trace(2, "DEC (%s+%d)", cpu.Reg.Name16(reg), d)
	return 23
}

func INC_IX_IY(cpu *CPU, reg *uint16) int {
	*reg++

	log.Trace(2, "INC %s", cpu.Reg.Name16(reg))
	return 10
}

func DEC_IX_IY(cpu *CPU, reg *uint16) int {
	*reg--

	log.Trace(2, "DEC %s", cpu.Reg.Name16(reg))
	return 10
}
