package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

func LD_R_R(cpu *CPU, dst *uint8, src *uint8) int {
	*dst = *src

	log.Trace(2, "LD %s, %s", cpu.Reg.Name(dst), cpu.Reg.Name(src))
	return 4
}

func UN_LD_R_R_HL(cpu *CPU, dst *uint8, src *uint16, low bool) int {
	l, h := helpers.To8(*src)
	if low {
		*dst = l
	} else {
		*dst = h
	}

	log.Trace(2, "LD %s, %s", cpu.Reg.Name(dst), cpu.Reg.Name16HL(src, low))
	return 4
}

func UN_LD_R_HL_R(cpu *CPU, dst *uint16, src *uint8, low bool) int {
	l, h := helpers.To8(*dst)
	if low {
		l = *src
	} else {
		h = *src
	}
	*dst = helpers.To16(l, h)

	log.Trace(2, "LD %s, %s", cpu.Reg.Name16HL(dst, low), cpu.Reg.Name(src))
	return 4
}

func LD_RR_nn(cpu *CPU, rh *uint8, rl *uint8) int {
	nn := FetchOperand16(cpu)

	*rl, *rh = helpers.To8(nn)

	log.Trace(2, "LD %s%s, $%04x", cpu.Reg.Name(rh), cpu.Reg.Name(rl), nn)
	return 10
}

func LD_mem_R_16(cpu *CPU, rh *uint8, rl *uint8, reg *uint8) int {
	MemoryWrite(cpu, helpers.To16(*rl, *rh), *reg)

	log.Trace(2, "LD (%s%s), %s", cpu.Reg.Name(rh), cpu.Reg.Name(rl), cpu.Reg.Name(reg))
	return 10
}

func LD_R_mem_16(cpu *CPU, rh *uint8, rl *uint8, reg *uint8) int {
	*reg = MemoryRead(cpu, helpers.To16(*rl, *rh))

	log.Trace(2, "LD %s, (%s%s)", cpu.Reg.Name(reg), cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 7
}

func LD_R_n(cpu *CPU, reg *uint8, n uint8) int {
	*reg = n

	log.Trace(2, "LD %s, $%02x", cpu.Reg.Name(reg), n)
	return 7
}

func POP(cpu *CPU, rh *uint8, rl *uint8) int {
	*rl, *rh = helpers.To8(PopStack16(cpu))

	log.Trace(2, "POP %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 10
}

func PUSH(cpu *CPU, rh *uint8, rl *uint8) int {
	PushStack16(cpu, helpers.To16(*rl, *rh))

	log.Trace(2, "PUSH %s%s", cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 11
}

func LD_R(cpu *CPU, reg *uint8) int {
	n := FetchOperand8(cpu)
	*reg = n
	log.Trace(2, "LD %s, $%02x", cpu.Reg.Name(reg), n)
	return 7
}

// IX / IY
func POP_IX_IY(cpu *CPU, reg *uint16) int {
	*reg = PopStack16(cpu)

	log.Trace(2, "POP %s", cpu.Reg.Name16(reg))
	return 15
}

func PUSH_IX_IY(cpu *CPU, reg *uint16) int {
	PushStack16(cpu, *reg)

	log.Trace(2, "PUSH %s", cpu.Reg.Name16(reg))
	return 15
}

func LD_R_IXIYd(cpu *CPU, idx *uint16, reg *uint8) int {
	d := FetchOperand8Compl(cpu)

	*reg = MemoryRead(cpu, uint16(int(*idx)+d))

	log.Trace(2, "LD %s, (%s+$%02x)", cpu.Reg.Name(reg), cpu.Reg.Name16(idx))
	return 19
}

func LD_IXIYd_R(cpu *CPU, idx *uint16, reg *uint8) int {
	d := FetchOperand8Compl(cpu)

	MemoryWrite(cpu, uint16(int(*idx)+d), *reg)

	log.Trace(2, "LD (%s+$%02x), %s", cpu.Reg.Name16(idx), d, cpu.Reg.Name(reg))
	return 19
}

func LD_IXIYd_n(cpu *CPU, idx *uint16) int {
	d := FetchOperand8Compl(cpu)
	n := FetchOperand8(cpu)

	MemoryWrite(cpu, uint16(int(*idx)+d), n)

	log.Trace(2, "LD (%s+$%02x), $%02x", cpu.Reg.Name16(idx), d, n)
	return 19
}

func LD_A_nn_mem(cpu *CPU) int {
	nn := FetchOperand16(cpu)
	cpu.Reg.A = MemoryRead(cpu, nn)

	log.Trace(2, "LD A, ($%04x)", nn)
	return 13
}

func LD_mem_R(cpu *CPU, reg *uint8) int {
	nn := FetchOperand16(cpu)

	MemoryWrite(cpu, nn, *reg)

	log.Trace(2, "LD (%04xh), %s", nn, cpu.Reg.Name(reg))
	return 13
}

func LD_A_I(cpu *CPU) int {
	cpu.Reg.A = cpu.Reg.I

	cpu.SetFlag(FLAG_SIGN, IsNeg(cpu.Reg.I))
	cpu.SetFlag(FLAG_ZERO, cpu.Reg.I == 0)
	cpu.SetFlag(FLAG_HALF_CARRY, false)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.IFF2)
	cpu.SetFlag(FLAG_ADD_SUB, false)

	// TODO: If an interrupt occurs during execution of this instruction, the Parity flag contains a 0
	log.Trace(2, "LD A, I")
	return 9
}

func LD_A_R(cpu *CPU) int {
	cpu.Reg.A = cpu.Reg.R

	cpu.SetFlag(FLAG_SIGN, IsNeg(cpu.Reg.R))
	cpu.SetFlag(FLAG_ZERO, cpu.Reg.R == 0)
	cpu.SetFlag(FLAG_HALF_CARRY, false)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.IFF2)
	cpu.SetFlag(FLAG_ADD_SUB, false)

	// TODO: If an interrupt occurs during execution of this instruction, the Parity flag contains a 0
	log.Trace(2, "LD A, R")
	return 9
}

func LD_I_A(cpu *CPU) int {
	cpu.Reg.I = cpu.Reg.A

	log.Trace(2, "LD I, A")
	return 9
}

func LD_R_A(cpu *CPU) int {
	cpu.Reg.R = cpu.Reg.A

	log.Trace(2, "LD R, A")
	return 9
}

func LD_R_nn(cpu *CPU, rh *uint8, rl *uint8) int {
	nn := FetchOperand16(cpu)

	*rl, *rh = helpers.To8(nn)

	log.Trace(2, "LD %s%s, $%04x", cpu.Reg.Name(rh), cpu.Reg.Name(rl), nn)
	return 10
}

func LD_IXIY_nn(cpu *CPU, reg *uint16) int {
	nn := FetchOperand16(cpu)

	*reg = nn

	log.Trace(2, "LD %s, $%04x", cpu.Reg.Name16(reg), nn)
	return 14
}

func LD_RR_nn_mem(cpu *CPU, rh *uint8, rl *uint8) int {
	nn := FetchOperand16(cpu)

	*rl, *rh = helpers.To8(MemoryRead16(cpu, nn))

	log.Trace(2, "LD %s%s, ($%04x)", cpu.Reg.Name(rh), cpu.Reg.Name(rl), nn)
	return 20
}

func LD_nn_RR_mem(cpu *CPU, rh *uint8, rl *uint8) int {
	nn := FetchOperand16(cpu)

	MemoryWrite16(cpu, nn, helpers.To16(*rl, *rh))

	log.Trace(2, "LD ($%04x), %s%s", nn, cpu.Reg.Name(rh), cpu.Reg.Name(rl))
	return 20
}

func LD_HL_nn_mem(cpu *CPU) int {
	nn := FetchOperand16(cpu)

	cpu.Reg.L = MemoryRead(cpu, nn)
	cpu.Reg.H = MemoryRead(cpu, nn+1)

	log.Trace(2, "LD HL, ($%04x)", nn)
	return 16
}

func LD_nn_HL_mem(cpu *CPU) int {
	nn := FetchOperand16(cpu)

	MemoryWrite16(cpu, nn, cpu.Reg.HL())

	log.Trace(2, "LD ($%04x), HL", nn)
	return 16
}

func LD_nn_IXIY_mem(cpu *CPU, reg *uint16) int {
	nn := FetchOperand16(cpu)

	MemoryWrite16(cpu, nn, *reg)

	log.Trace(2, "LD ($%04x), %s", nn, cpu.Reg.Name16(reg))
	return 20
}

func LD_IXIY_nn_mem(cpu *CPU, reg *uint16) int {
	nn := FetchOperand16(cpu)

	*reg = MemoryRead16(cpu, nn)

	log.Trace(2, "LD %s, ($%04x)", cpu.Reg.Name16(reg), nn)
	return 20
}

func LD_SP_IX_IY(cpu *CPU, idx *uint16) int {
	cpu.Reg.SP = *idx

	log.Trace(2, "LD SP, %s", cpu.Reg.Name16(idx))
	return 10
}
