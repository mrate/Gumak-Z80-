package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

func IN_A_n(cpu *CPU) int {
	al := FetchOperand8(cpu)
	ah := cpu.Reg.A

	cpu.Reg.A = ReadIO(cpu, helpers.To16(al, ah))

	log.Trace(2, "IN A, ($%02x)", al)
	return 11
}

func IN_R_C(cpu *CPU, reg *uint8) int {
	data := ReadIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B))
	*reg = data

	cpu.SetFlag(FLAG_SIGN, (data&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_ZERO, data == 0)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_P_V, Alu_ParityEven(data))
	cpu.SetFlag(FLAG_N, false)

	log.Trace(2, "IN %s, (C)", cpu.Reg.Name(reg))
	return 12
}

func INI(cpu *CPU) int {
	data := ReadIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B))

	MEM_HL_W(cpu, data)

	cpu.Reg.B--
	cpu.Reg.HL_write(cpu.Reg.HL() + 1)

	cpu.SetFlag(FLAG_ZERO, cpu.Reg.B == 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	log.Trace(2, "INI")
	return 16
}

func INIR(cpu *CPU) int {
	INI(cpu)
	log.Trace(2, "R")

	if cpu.Reg.B != 0 {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}

func IND(cpu *CPU) int {
	data := ReadIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B))

	MEM_HL_W(cpu, data)

	cpu.Reg.B--
	cpu.Reg.HL_write(cpu.Reg.HL() - 1)

	cpu.SetFlag(FLAG_ZERO, cpu.Reg.B == 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	log.Trace(2, "IND")
	return 16
}

func INDR(cpu *CPU) int {
	IND(cpu)
	log.Trace(2, "R")

	if cpu.Reg.B != 0 {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}

func OUT_n_A(cpu *CPU) int {
	al := FetchOperand8(cpu)
	ah := cpu.Reg.A

	WriteIO(cpu, helpers.To16(al, ah), cpu.Reg.A)

	log.Trace(2, "OUT ($%02x), A", al)
	return 12
}

func OUTI(cpu *CPU) int {
	data := MEM_HL(cpu)
	cpu.Reg.B--
	WriteIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B), data)
	cpu.Reg.HL_write(cpu.Reg.HL() + 1)

	cpu.SetFlag(FLAG_ZERO, cpu.Reg.B == 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	log.Trace(2, "OUTI")
	return 16
}

func OUT_C_R(cpu *CPU, reg *uint8) int {
	WriteIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B), *reg)

	log.Trace(2, "OUT (C), %s", cpu.Reg.Name(reg))
	return 12
}

func OUTIR(cpu *CPU) int {
	OUTI(cpu)
	log.Trace(2, "R")

	if cpu.Reg.B != 0 {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}

func OUTD(cpu *CPU) int {
	data := MEM_HL(cpu)
	cpu.Reg.B--
	WriteIO(cpu, helpers.To16(cpu.Reg.C, cpu.Reg.B), data)
	cpu.Reg.HL_write(cpu.Reg.HL() - 1)

	cpu.SetFlag(FLAG_ZERO, cpu.Reg.B == 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	log.Trace(2, "OUTD")
	return 16
}

func OTDR(cpu *CPU) int {
	OUTD(cpu)
	log.Trace(2, "R")

	if cpu.Reg.B != 0 {
		cpu.Reg.PC -= 2
		return 21
	}

	return 16
}
