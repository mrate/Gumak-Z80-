package z80

import (
	"mutex/gumak/log"
)

func BitStr(flags uint8, bit uint8, c string) string {
	if (flags & bit) != 0 {
		return c
	} else {
		return "-"
	}
}

func FlagsToStr(flags uint8) string {
	res := ""

	res += BitStr(flags, FLAG_CARRY, "C")
	res += BitStr(flags, FLAG_ADD_SUB, "N")
	res += BitStr(flags, FLAG_PARTY_OVERFLOW, "P")
	res += BitStr(flags, FLAG_HALF_CARRY, "H")
	res += BitStr(flags, FLAG_ZERO, "Z")
	res += BitStr(flags, FLAG_SIGN, "S")

	return res
}

func (cpu *CPU) DumpState() {
	log.Debug("=====================================")
	log.Debug("CPU [%d] [%0.2f us]", cpu.Frequency, cpu.TStateUs)
	log.Debug("")
	log.Debug("PC: 0x%04x [%s]  SP: 0x%04x  IFF1: %t  IFF2: %t", cpu.Reg.PC, SymbolForAddressRelative(cpu, cpu.Reg.PC), cpu.Reg.SP, cpu.IFF1, cpu.IFF2)
	log.Debug("Halted: %t", cpu.halted)
	log.Debug("")
	log.Debug("A: 0x%02x  F: 0x%02x  [%s]", cpu.Reg.A, cpu.Reg.F, FlagsToStr(cpu.Reg.F))
	log.Debug("B: 0x%02x  C: 0x%02x", cpu.Reg.B, cpu.Reg.C)
	log.Debug("D: 0x%02x  E: 0x%02x", cpu.Reg.D, cpu.Reg.E)
	log.Debug("H: 0x%02x  L: 0x%02x", cpu.Reg.H, cpu.Reg.L)
	log.Debug("")
	log.Debug("A': 0x%02x  F': 0x%02x  [%s]", cpu.Reg.A_, cpu.Reg.F_, FlagsToStr(cpu.Reg.F_))
	log.Debug("B': 0x%02x  C': 0x%02x", cpu.Reg.B_, cpu.Reg.C_)
	log.Debug("D': 0x%02x  E': 0x%02x", cpu.Reg.D_, cpu.Reg.E_)
	log.Debug("H': 0x%02x  L': 0x%02x", cpu.Reg.H_, cpu.Reg.L_)
	log.Debug("")
	log.Debug("I/O: Addr: 0x%04x  Data: 0x%02x", cpu.Pin.ADDR, cpu.Pin.DATA)
	log.Debug("=====================================")
}
