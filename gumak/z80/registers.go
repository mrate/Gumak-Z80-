package z80

import "mutex/gumak/helpers"

type Registers struct {
	A  uint8  // Accumulator
	F  uint8  // Flag
	B  uint8  // Byte count
	C  uint8  //
	D  uint8  //
	E  uint8  //
	H  uint8  //
	L  uint8  //
	A_ uint8  // Accumulator'
	F_ uint8  // Flag'
	B_ uint8  // B'
	C_ uint8  // C'
	D_ uint8  // D'
	E_ uint8  // E'
	H_ uint8  // H'
	L_ uint8  // L'
	IX uint16 // Index register 1
	IY uint16 // Index register 2
	SP uint16 // Stack pointer
	PC uint16 // Program counter
	R  uint8  // Memory refresh
	R7 uint8  // Memory refresh
	I  uint8  // Interrupt
}

// Bit:   7   6   5   4   3   2   1   0
// Pos:   S   Z   X   H   X  P/V  N   C

const (
	FLAG_SIGN           uint8 = 1 << 7 // S
	FLAG_ZERO           uint8 = 1 << 6 // Z
	FLAG_5              uint8 = 1 << 5 // Unused
	FLAG_HALF_CARRY     uint8 = 1 << 4 // H
	FLAG_3              uint8 = 1 << 3 // Unused
	FLAG_PARTY_OVERFLOW uint8 = 1 << 2 // P/V
	FLAG_ADD_SUB        uint8 = 1 << 1 // N
	FLAG_CARRY          uint8 = 1 << 0 // C

	// Aliases.
	FLAG_S   uint8 = FLAG_SIGN
	FLAG_Z   uint8 = FLAG_ZERO
	FLAG_H   uint8 = FLAG_HALF_CARRY
	FLAG_P_V uint8 = FLAG_PARTY_OVERFLOW
	FLAG_N   uint8 = FLAG_ADD_SUB
	FLAG_C   uint8 = FLAG_CARRY
)

// Helpers.
func (r *Registers) Name(reg *uint8) string {
	switch reg {
	case &r.A:
		return "A"
	case &r.F:
		return "F"
	case &r.B:
		return "B"
	case &r.C:
		return "C"
	case &r.D:
		return "D"
	case &r.E:
		return "E"
	case &r.H:
		return "H"
	case &r.L:
		return "L"
	case &r.A_:
		return "A'"
	case &r.F_:
		return "F'"
	case &r.B_:
		return "B'"
	case &r.C_:
		return "C'"
	case &r.D_:
		return "D'"
	case &r.E_:
		return "E'"
	case &r.H_:
		return "H'"
	case &r.L_:
		return "L'"
	case &r.R:
		return "R"
	case &r.R7:
		return "R7"
	case &r.I:
		return "I"
	default:
		return "?"
	}
}

func (r *Registers) Name16(reg *uint16) string {
	switch reg {
	case &r.IX:
		return "IX"
	case &r.IY:
		return "IY"
	case &r.SP:
		return "SP"
	case &r.PC:
		return "PC"
	default:
		return "?"
	}
}

func (r *Registers) Name16HL(reg *uint16, low bool) string {
	var suffix string
	if low {
		suffix = "l"
	} else {
		suffix = "h"
	}

	return r.Name16(reg) + suffix
}

func (r *Registers) AF() uint16 {
	return helpers.To16(r.F, r.A)
}

func (r *Registers) BC() uint16 {
	return helpers.To16(r.C, r.B)
}

func (r *Registers) DE() uint16 {
	return helpers.To16(r.E, r.D)
}

func (r *Registers) HL() uint16 {
	return helpers.To16(r.L, r.H)
}

func (r *Registers) AF_write(value uint16) {
	r.F, r.A = helpers.To8(value)
}

func (r *Registers) BC_write(value uint16) {
	r.C, r.B = helpers.To8(value)
}

func (r *Registers) DE_write(value uint16) {
	r.E, r.D = helpers.To8(value)
}

func (r *Registers) HL_write(value uint16) {
	r.L, r.H = helpers.To8(value)
}

func (r *Registers) R_() uint8 {
	return (r.R & 0x7f) | r.R7
}

func (r *Registers) R_write(value uint8) {
	r.R = value
	r.R7 = value & 0x80
}

func (r *Registers) Clear() {
	r.A = 0
	r.B = 0
	r.C = 0
	r.D = 0
	r.E = 0
	r.F = 0
	r.H = 0
	r.L = 0
	r.A_ = 0
	r.B_ = 0
	r.C_ = 0
	r.D_ = 0
	r.E_ = 0
	r.F_ = 0
	r.H_ = 0
	r.L_ = 0

	r.IX = 0
	r.IY = 0
	r.SP = 0
	r.PC = 0
	r.R = 0
	r.R7 = 0
	r.I = 0
}

func (cpu *CPU) SetFlag(bit uint8, set bool) {
	if set {
		cpu.Reg.F = cpu.Reg.F | bit
	} else {
		cpu.Reg.F = cpu.Reg.F &^ bit
	}
}

func (cpu *CPU) Flag(bit uint8) bool {
	return (cpu.Reg.F & bit) != 0
}

func IsNeg(f uint8) bool {
	return f&FLAG_SIGN != 0
}

func FlagName(flag uint8, value bool) string {
	if value {
		switch flag {
		case FLAG_SIGN:
			return "S"
		case FLAG_ZERO:
			return "Z"
		case FLAG_HALF_CARRY:
			return "H"
		case FLAG_PARTY_OVERFLOW:
			return "P/V"
		case FLAG_ADD_SUB:
			return "N"
		case FLAG_CARRY:
			return "C"
		case FLAG_3:
			return "3"
		case FLAG_5:
			return "5"
		default:
			return "??"
		}
	} else {
		switch flag {
		case FLAG_SIGN:
			return "NS"
		case FLAG_ZERO:
			return "NZ"
		case FLAG_HALF_CARRY:
			return "NH"
		case FLAG_PARTY_OVERFLOW:
			return "NP/V"
		case FLAG_ADD_SUB:
			return "NN"
		case FLAG_CARRY:
			return "NC"
		case FLAG_3:
			return "N3"
		case FLAG_5:
			return "N5"
		default:
			return "N??"
		}
	}
}
