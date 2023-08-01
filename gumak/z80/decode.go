package z80

import (
	"mutex/gumak/helpers"
	"mutex/gumak/log"
)

func DecodeInstruction(cpu *CPU) InstrOp {
	op := FetchInstruction(cpu)
	return OpCodes[op]
}

func DecodeAndExecute(cpu *CPU) int {
	if cpu.symbols != nil {
		log.Trace(2, "[%s] ", SymbolForAddressRelative(cpu, cpu.Reg.PC))
	}

	if cpu.breakPoints != nil {
		if cb, ok := cpu.breakPoints[cpu.Reg.PC]; ok {
			cb()
		}
	}

	log.Trace(2, "[PC: 0x%04x] ", cpu.Reg.PC)

	instr := DecodeInstruction(cpu)
	tStates := instr(cpu)

	log.TraceFlush()

	return tStates
}

func FetchInstruction(cpu *CPU) uint8 {
	inst := MemoryRead(cpu, cpu.Reg.PC)
	cpu.Reg.PC++
	return inst
}

func FetchOperand8(cpu *CPU) uint8 {
	return FetchInstruction(cpu)
}

func FetchOperand16(cpu *CPU) uint16 {
	l := FetchOperand8(cpu)
	h := FetchOperand8(cpu)

	return helpers.To16(l, h)
}

func FetchOperand8Compl(cpu *CPU) int {
	return int(int8(FetchOperand8(cpu)))
}
