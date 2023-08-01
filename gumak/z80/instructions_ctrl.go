package z80

import "mutex/gumak/log"

func DAA(cpu *CPU) int {
	log.Trace(2, "DAA")

	a := cpu.Reg.A
	incr := uint8(0)
	carry := cpu.Flag(FLAG_CARRY)

	if cpu.Flag(FLAG_HALF_CARRY) || (a&0x0f) > 0x09 {
		incr |= 0x06
	}

	if (carry || (a > 0x9f)) || ((a > 0x8f) && ((a & 0x0f) > 0x09)) {
		incr |= 0x60
	}

	if a > 0x99 {
		carry = true
	}

	if cpu.Flag(FLAG_ADD_SUB) {
		Alu_SUB_A(cpu, incr)
	} else {
		Alu_ADD_A(cpu, incr)
	}

	cpu.SetFlag(FLAG_PARTY_OVERFLOW, Alu_ParityEven(cpu.Reg.A))
	cpu.SetFlag(FLAG_CARRY, carry)

	return 4
}

func HALT(cpu *CPU) int {
	log.Trace(2, "HALT")

	cpu.halted = true
	return 4
}

func DI(cpu *CPU) int {
	log.Trace(2, "DI")

	cpu.IFF1 = false
	cpu.IFF2 = false
	cpu.maskableSkip = 0
	return 4
}

func EI(cpu *CPU) int {
	log.Trace(2, "EI")

	cpu.IFF1 = true
	cpu.IFF2 = true
	cpu.maskableSkip = 1
	return 4
}

func IM0(cpu *CPU) int {
	log.Trace(2, "IM 0")

	cpu.InterruptMode = 0
	return 8
}

func IM1(cpu *CPU) int {
	log.Trace(2, "IM 1")

	cpu.InterruptMode = 1
	return 8
}

func IM2(cpu *CPU) int {
	log.Trace(2, "IM 2")

	cpu.InterruptMode = 2
	return 8
}
