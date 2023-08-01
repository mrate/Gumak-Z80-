package z80

import "mutex/gumak/log"

func JP_nn(cpu *CPU) int {
	// OK
	nn := FetchOperand16(cpu)
	cpu.Reg.PC = nn

	log.Trace(2, "JP $%04x", nn)
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}

	return 10
}

func JP_FLAG_nn(cpu *CPU, flag uint8, value bool) int {
	// OK
	nn := FetchOperand16(cpu)

	if cpu.Flag(flag) == value {
		cpu.Reg.PC = nn
	}

	log.Trace(2, "JP %s, $%04x", FlagName(flag, value), nn)
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, nn))
	}

	return 10
}

func JR_e(cpu *CPU) int {
	// OK
	e := FetchOperand8Compl(cpu)

	cpu.Reg.PC = uint16(int32(cpu.Reg.PC) + int32(e))

	log.Trace(2, "JR %+d", e)
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}

	return 12
}

func JR_FLAG_e(cpu *CPU, flag uint8, value bool) int {
	// OK
	e := FetchOperand8Compl(cpu)

	log.Trace(2, "JR [%s], %+d", FlagName(flag, value), e)
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, uint16(int32(cpu.Reg.PC)+int32(e))))
	}

	if cpu.Flag(flag) == value {
		cpu.Reg.PC = uint16(int32(cpu.Reg.PC) + int32(e))
		return 12
	}

	return 7
}

func JP_IHLI(cpu *CPU) int {
	// OK
	cpu.Reg.PC = cpu.Reg.HL()

	log.Trace(2, "JP (HL)")
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}
	return 4
}

func JP_IX_IY(cpu *CPU, idx *uint16) int {
	// OK
	cpu.Reg.PC = *idx

	log.Trace(2, "JP %s", cpu.Reg.Name16(idx))
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}
	return 8
}

func DJNZ_e(cpu *CPU) int {
	// OK
	e := FetchOperand8Compl(cpu)
	cpu.Reg.B--

	log.Trace(2, "DJNZ %+d", e)
	if cpu.symbols != nil {
		log.Trace(2, " [%s]", SymbolForAddressRelative(cpu, uint16(int32(cpu.Reg.PC)+int32(e))))
	}

	if cpu.Reg.B != 0 {
		cpu.Reg.PC = uint16(int32(cpu.Reg.PC) + int32(e))
		return 3
	}

	return 2
}
