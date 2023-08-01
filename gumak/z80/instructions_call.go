package z80

import "mutex/gumak/log"

func CALL_nn(cpu *CPU) int {
	pc := FetchOperand16(cpu)

	PushStack16(cpu, cpu.Reg.PC)
	cpu.Reg.PC = pc

	log.Trace(1, "CALL $%04x", cpu.Reg.PC)
	if cpu.symbols != nil {
		log.Trace(1, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}

	return 17
}

func CALL_FLAG_nn(cpu *CPU, flag uint8, value bool) int {
	pc := FetchOperand16(cpu)

	log.Trace(1, "CALL %s, $%04x", FlagName(flag, value), pc)
	if cpu.symbols != nil {
		log.Trace(1, " [%s]", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}

	if cpu.Flag(flag) == value {
		PushStack16(cpu, cpu.Reg.PC)
		cpu.Reg.PC = pc

		return 17
	}

	return 10
}

func RET(cpu *CPU) int {
	cpu.Reg.PC = PopStack16(cpu)

	log.Trace(2, "RET")
	return 10
}

func RET_cc(cpu *CPU, flag uint8, value bool) int {
	log.Trace(2, "RET %s", FlagName(flag, value))

	if cpu.Flag(flag) == value {
		cpu.Reg.PC = PopStack16(cpu)

		return 11
	}

	return 5
}

func RETI(cpu *CPU) int {
	cpu.Reg.PC = PopStack16(cpu)

	// TODO: Signal I/O device

	log.Trace(2, "RETI")
	return 14
}

func RETN(cpu *CPU) int {
	cpu.Reg.PC = PopStack16(cpu)
	cpu.IFF1 = cpu.IFF2

	// TODO: Signal I/O device

	log.Trace(2, "RETN")
	return 14
}

func RST(cpu *CPU, val uint8) int {
	PushStack16(cpu, cpu.Reg.PC)
	cpu.Reg.PC = uint16(val)

	log.Trace(1, "RST $%02x", val)

	return 11
}
