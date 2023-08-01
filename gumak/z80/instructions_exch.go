package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

func EX_DE_HL(cpu *CPU) int {
	de := cpu.Reg.DE()
	hl := cpu.Reg.HL()

	cpu.Reg.DE_write(hl)
	cpu.Reg.HL_write(de)

	log.Trace(2, "EX DE, HL")
	return 4
}

func EX_AF_AF_(cpu *CPU) int {
	helpers.Exchange8(&cpu.Reg.A, &cpu.Reg.A_)
	helpers.Exchange8(&cpu.Reg.F, &cpu.Reg.F_)

	log.Trace(2, "EX AF, AF'")
	return 4
}

func EXX(cpu *CPU) int {
	// OK
	helpers.Exchange8(&cpu.Reg.B, &cpu.Reg.B_)
	helpers.Exchange8(&cpu.Reg.C, &cpu.Reg.C_)
	helpers.Exchange8(&cpu.Reg.D, &cpu.Reg.D_)
	helpers.Exchange8(&cpu.Reg.E, &cpu.Reg.E_)
	helpers.Exchange8(&cpu.Reg.H, &cpu.Reg.H_)
	helpers.Exchange8(&cpu.Reg.L, &cpu.Reg.L_)

	log.Trace(2, "EXX")
	return 4
}

func EX_ISPI_HL(cpu *CPU) int {
	// OK
	h := cpu.Reg.H
	l := cpu.Reg.L

	cpu.Reg.L = MEM_SP(cpu, 0)
	cpu.Reg.H = MEM_SP(cpu, 1)

	MEM_SP_W(cpu, 0, l)
	MEM_SP_W(cpu, 1, h)

	log.Trace(2, "EX (SP), HL")
	return 19
}

func EX_SP_IX_IY(cpu *CPU, idx *uint16) int {
	// OK
	l, h := helpers.To8(*idx)

	*idx = helpers.To16(MEM_SP(cpu, 0), MEM_SP(cpu, 1))
	MEM_SP_W(cpu, 0, l)
	MEM_SP_W(cpu, 1, h)

	log.Trace(2, "EX (SP), %s", cpu.Reg.Name16(idx))
	return 23
}

func LDI(cpu *CPU) int {
	// OK
	MEM_DE_W(cpu, MEM_HL(cpu))

	cpu.Reg.DE_write(cpu.Reg.DE() + 1)
	cpu.Reg.HL_write(cpu.Reg.HL() + 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.Reg.BC() != 0)

	log.Trace(2, "LDI")
	return 16
}

func LDIR(cpu *CPU) int {
	// OK
	LDI(cpu)
	cpu.Refresh(2)
	log.Trace(2, "R")

	if cpu.Reg.BC() != 0 {
		cpu.Reg.PC -= 2
		return 23
	}

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, false)
	return 16
}

func LDD(cpu *CPU) int {
	// OK
	MEM_DE_W(cpu, MEM_HL(cpu))

	cpu.Reg.DE_write(cpu.Reg.DE() - 1)
	cpu.Reg.HL_write(cpu.Reg.HL() - 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.Reg.BC() != 0)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)

	log.Trace(2, "LDD")
	return 16
}

func LDDR(cpu *CPU) int {
	// OK
	LDD(cpu)
	cpu.Refresh(2)
	log.Trace(2, "R")

	if cpu.Reg.BC() != 0 {
		cpu.Reg.PC -= 2
		return 23
	}

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, false)
	return 16
}

func CPI(cpu *CPU) int {
	// OK
	c := cpu.Flag(FLAG_CARRY)

	Alu_CP_A(cpu, MEM_HL(cpu))
	cpu.Reg.HL_write(cpu.Reg.HL() + 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.Reg.BC() != 0)
	cpu.SetFlag(FLAG_CARRY, c)

	log.Trace(2, "CPI")
	return 16
}

func CPIR(cpu *CPU) int {
	// OK
	c := cpu.Flag(FLAG_CARRY)

	Alu_CP_A(cpu, MEM_HL(cpu))
	cpu.Reg.HL_write(cpu.Reg.HL() + 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	pv := cpu.Reg.BC() != 0

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, pv)
	cpu.SetFlag(FLAG_CARRY, c)

	log.Trace(2, "CPIR")
	if pv && !cpu.Flag(FLAG_ZERO) {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}

func CPD(cpu *CPU) int {
	// OK
	c := cpu.Flag(FLAG_CARRY)

	Alu_CP_A(cpu, MEM_HL(cpu))
	cpu.Reg.HL_write(cpu.Reg.HL() - 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, cpu.Reg.BC() != 0)
	cpu.SetFlag(FLAG_CARRY, c)

	log.Trace(2, "CPD")
	return 16
}

func CPDR(cpu *CPU) int {
	// OK
	c := cpu.Flag(FLAG_CARRY)

	Alu_CP_A(cpu, MEM_HL(cpu))
	cpu.Reg.HL_write(cpu.Reg.HL() - 1)
	cpu.Reg.BC_write(cpu.Reg.BC() - 1)

	pv := cpu.Reg.BC() != 0

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, pv)
	cpu.SetFlag(FLAG_CARRY, c)

	log.Trace(2, "CPDR")
	if pv && !cpu.Flag(FLAG_ZERO) {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}
