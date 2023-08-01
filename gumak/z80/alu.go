package z80

import "mutex/gumak/helpers"

var parity [256]bool

func init() {
	for i := range parity {
		parity[i] = helpers.CountBits(i)%2 == 0
	}
}

func Alu_ParityEven(value uint8) bool {
	return parity[value]
}

func Alu_ADC_A(cpu *CPU, b uint8) {
	a := cpu.Reg.A
	c := uint8(0)
	if cpu.Flag(FLAG_CARRY) {
		c = 1
	}
	wres := uint16(a) + uint16(b) + uint16(c)
	res := uint8(wres & 0xff)

	// TODO:
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)

	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_CARRY, (wres&0x100) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, ((a^(^b))&(a^res)&0x80) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, (((a&0x0f)+(b&0x0f)+c)&FLAG_HALF_CARRY) != 0)
	cpu.SetFlag(FLAG_ADD_SUB, false)

	cpu.Reg.A = res
}

func Alu_ADD_A(cpu *CPU, b uint8) {
	a := cpu.Reg.A
	wres := uint16(a) + uint16(b)
	res := uint8(wres & 0xff)

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_CARRY, (wres&0x100) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, ((a^(^b))&(a^res)&0x80) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, (((a&0x0f)+(b&0x0f))&FLAG_HALF_CARRY) != 0)
	cpu.SetFlag(FLAG_ADD_SUB, false)

	cpu.Reg.A = res
}

func Alu_SBC_A(cpu *CPU, b uint8) {
	a := cpu.Reg.A
	c := uint8(0)
	if cpu.Flag(FLAG_CARRY) {
		c = 1
	}
	wres := uint16(a) - uint16(b) - uint16(c)
	res := uint8(wres & 0xff)

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_CARRY, (wres&0x100) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, ((a^(^b))&(a^res)&0x80) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, (((a&0x0f)-(b&0x0f)-c)&FLAG_HALF_CARRY) != 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	cpu.Reg.A = res
}

func Alu_SUB_A(cpu *CPU, b uint8) {
	a := cpu.Reg.A
	wres := uint16(a) - uint16(b)
	res := uint8(wres & 0xff)

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_CARRY, (wres&0x100) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, ((a^(^b))&(a^res)&0x80) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, (((a&0x0f)-(b&0x0f))&FLAG_HALF_CARRY) != 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)

	cpu.Reg.A = res
}

func Alu_AND_A(cpu *CPU, b uint8) {
	res := cpu.Reg.A & b

	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, true)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, (parity[res]))
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_ADD_SUB, false)
	cpu.SetFlag(FLAG_CARRY, false)

	cpu.Reg.A = res
}

func Alu_OR_A(cpu *CPU, b uint8) {
	res := cpu.Reg.A | b

	cpu.SetFlag(FLAG_S, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, false)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, parity[res])
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_ADD_SUB, false)
	cpu.SetFlag(FLAG_CARRY, false)

	cpu.Reg.A = res
}

func Alu_XOR_A(cpu *CPU, b uint8) {
	res := (cpu.Reg.A ^ b) & 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_Z, res == 0)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, false)

	cpu.Reg.A = res
}

func Alu_NEG_A(cpu *CPU) {
	a := cpu.Reg.A

	prev := a
	a = 0 - a

	cpu.Reg.A = a

	cpu.SetFlag(FLAG_SIGN, IsNeg(a))
	cpu.SetFlag(FLAG_ZERO, a == 0)
	cpu.SetFlag(FLAG_HALF_CARRY, false)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, prev == 0x80)
	cpu.SetFlag(FLAG_ADD_SUB, true)
	cpu.SetFlag(FLAG_CARRY, prev != 0x0)
}

func Alu_CP_A(cpu *CPU, b uint8) {
	a := cpu.Reg.A
	wres := uint16(a) - uint16(b)
	res := uint8(wres & 0xff)

	cpu.SetFlag(FLAG_SIGN, (res&FLAG_SIGN) != 0)
	cpu.SetFlag(FLAG_3, (b&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (b&FLAG_5) != 0)
	cpu.SetFlag(FLAG_ADD_SUB, true)
	cpu.SetFlag(FLAG_ZERO, res == 0)
	cpu.SetFlag(FLAG_CARRY, (wres&0x100) != 0)
	cpu.SetFlag(FLAG_HALF_CARRY, (((a&0x0f)-(b&0x0f))&FLAG_HALF_CARRY) != 0)
	cpu.SetFlag(FLAG_PARTY_OVERFLOW, ((a^b)&(a^res)&0x80) != 0)
}

func Alu_RLC_A(cpu *CPU) {
	res := cpu.Reg.A
	c := (res & 0x80) != 0

	if c {
		res = (res << 1) | 0x01
	} else {
		res <<= 1
	}
	res &= 0xff

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_C, c)

	cpu.Reg.A = res
}

func Alu_RRC_A(cpu *CPU) {
	res := cpu.Reg.A
	c := (res & 0x01) != 0

	if c {
		res = (res >> 1) | 0x80
	} else {
		res >>= 1
	}

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_C, c)

	cpu.Reg.A = res
}

func Alu_RLC(cpu *CPU, res uint8) uint8 {
	c := (res & 0x80) != 0

	if c {
		res = (res << 1) | 0x01
	} else {
		res <<= 1
	}
	res &= 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_RRC(cpu *CPU, res uint8) uint8 {
	c := (res & 0x01) != 0

	if c {
		res = (res >> 1) | 0x80
	} else {
		res >>= 1
	}

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_RL_A(cpu *CPU) {
	res := cpu.Reg.A
	c := (res & 0x80) != 0

	if cpu.Flag(FLAG_C) {
		res = (res << 1) | 0x01
	} else {
		res <<= 1
	}

	res &= 0xff

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_C, c)

	cpu.Reg.A = res
}

func Alu_RR_A(cpu *CPU) {
	res := cpu.Reg.A
	c := (res & 0x01) != 0

	if cpu.Flag(FLAG_C) {
		res = (res >> 1) | 0x80
	} else {
		res >>= 1
	}

	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_C, c)

	cpu.Reg.A = res
}

func Alu_CPL_A(cpu *CPU) {
	cpu.Reg.A = ^cpu.Reg.A

	cpu.SetFlag(FLAG_HALF_CARRY, true)
	cpu.SetFlag(FLAG_ADD_SUB, true)
}

func Alu_CCF(cpu *CPU) {
	cpu.SetFlag(FLAG_HALF_CARRY, cpu.Flag(FLAG_CARRY))
	cpu.SetFlag(FLAG_ADD_SUB, false)
	cpu.SetFlag(FLAG_CARRY, !cpu.Flag(FLAG_CARRY))
}

func Alu_SCF(cpu *CPU) {
	cpu.SetFlag(FLAG_HALF_CARRY, false)
	cpu.SetFlag(FLAG_ADD_SUB, false)
	cpu.SetFlag(FLAG_CARRY, true)
}

func Alu_RL(cpu *CPU, res uint8) uint8 {
	c := (res & 0x80) != 0

	if cpu.Flag(FLAG_C) {
		res = (res << 1) | 0x01
	} else {
		res <<= 1
	}
	res &= 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_RR(cpu *CPU, res uint8) uint8 {
	c := (res & 0x01) != 0

	if cpu.Flag(FLAG_C) {
		res = (res >> 1) | 0x80
	} else {
		res >>= 1
	}

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_SLA(cpu *CPU, res uint8) uint8 {
	c := (res & 0x80) != 0
	res = (res << 1) & 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_SLS(cpu *CPU, res uint8) uint8 {
	c := (res & 0x80) != 0
	res = ((res << 1) | 0x01) & 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_SRA(cpu *CPU, res uint8) uint8 {
	c := (res & 0x01) != 0
	res = (res >> 1) | (res & 0x80)

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_SRL(cpu *CPU, res uint8) uint8 {
	c := (res & 0x01) != 0
	res >>= 1

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, parity[res])
	cpu.SetFlag(FLAG_H, false)
	cpu.SetFlag(FLAG_N, false)
	cpu.SetFlag(FLAG_C, c)

	return res
}

func Alu_DEC8(cpu *CPU, res uint8) uint8 {
	pv := (res == 0x80)
	h := (((res & 0x0f) - 1) & FLAG_H) != 0
	res = (res - 1) & 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, pv)
	cpu.SetFlag(FLAG_H, h)
	cpu.SetFlag(FLAG_N, true)

	return res
}

func Alu_INC8(cpu *CPU, res uint8) uint8 {
	pv := (res == 0x7f)
	h := (((res & 0x0f) + 1) & FLAG_H) != 0
	res = (res + 1) & 0xff

	cpu.SetFlag(FLAG_S, (res&FLAG_S) != 0)
	cpu.SetFlag(FLAG_3, (res&FLAG_3) != 0)
	cpu.SetFlag(FLAG_5, (res&FLAG_5) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_P_V, pv)
	cpu.SetFlag(FLAG_H, h)
	cpu.SetFlag(FLAG_N, false)

	return res
}

func Alu_ADD16(cpu *CPU, a uint16, b uint16) uint16 {
	var lres int = int(a) + int(b)
	res := uint16(lres & 0xffff)

	cpu.SetFlag(FLAG_3, (res&(uint16(FLAG_3)<<8)) != 0)
	cpu.SetFlag(FLAG_5, (res&(uint16(FLAG_5)<<8)) != 0)
	cpu.SetFlag(FLAG_C, (lres&0x10000) != 0)
	cpu.SetFlag(FLAG_H, (((a&0x0fff)+(b&0x0fff))&0x1000) != 0)
	cpu.SetFlag(FLAG_N, false)

	return res
}

func Alu_ADC16(cpu *CPU, a uint16, b uint16) uint16 {
	c := uint16(0)
	if cpu.Flag(FLAG_CARRY) {
		c = 1
	}

	lres := int(a) + int(b) + int(c)
	res := uint16(lres & 0xffff)

	cpu.SetFlag(FLAG_S, (res&(uint16(FLAG_S)<<8)) != 0)
	cpu.SetFlag(FLAG_3, (res&(uint16(FLAG_3)<<8)) != 0)
	cpu.SetFlag(FLAG_5, (res&(uint16(FLAG_5)<<8)) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_C, (lres&0x10000) != 0)
	cpu.SetFlag(FLAG_P_V, ((a^(^b))&(a^res)&0x8000) != 0)
	cpu.SetFlag(FLAG_H, (((a&0x0fff)+(b&0x0fff)+c)&0x1000) != 0)
	cpu.SetFlag(FLAG_N, false)

	return res
}

func Alu_SBC16(cpu *CPU, a uint16, b uint16) uint16 {
	c := uint16(0)
	if cpu.Flag(FLAG_CARRY) {
		c = 1
	}

	var lres int = int(a) - int(b) - int(c)
	res := uint16(lres & 0xffff)

	cpu.SetFlag(FLAG_S, (res&(uint16(FLAG_S)<<8)) != 0)
	cpu.SetFlag(FLAG_3, (res&(uint16(FLAG_3)<<8)) != 0)
	cpu.SetFlag(FLAG_5, (res&(uint16(FLAG_5)<<8)) != 0)
	cpu.SetFlag(FLAG_Z, (res) == 0)
	cpu.SetFlag(FLAG_C, (lres&0x10000) != 0)
	cpu.SetFlag(FLAG_P_V, ((a^b)&(a^res)&0x8000) != 0)
	cpu.SetFlag(FLAG_H, (((a&0x0fff)-(b&0x0fff)-c)&0x1000) != 0)
	cpu.SetFlag(FLAG_N, true)

	return res
}

func Alu_BIT(cpu *CPU, res uint8, bit uint8) {
	cpu.SetFlag(FLAG_ZERO, (res&(1<<bit)) == 0)
	cpu.SetFlag(FLAG_HALF_CARRY, true)
	cpu.SetFlag(FLAG_N, false)
}

func Alu_RES(cpu *CPU, res uint8, bit uint8) uint8 {
	return res &^ (1 << bit)
}

func Alu_SET(cpu *CPU, res uint8, bit uint8) uint8 {
	return res | (1 << bit)
}
