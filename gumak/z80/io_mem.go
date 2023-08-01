package z80

import "mutex/gumak/helpers"

func MEM_HL(cpu *CPU) uint8 {
	return MemoryRead(cpu, cpu.Reg.HL())
}

func MEM_HL_W(cpu *CPU, value uint8) {
	MemoryWrite(cpu, cpu.Reg.HL(), value)
}

func MEM_IX(cpu *CPU, d int) uint8 {
	return MemoryRead(cpu, uint16(int(cpu.Reg.IX)+d))
}

func MEM_IX_W(cpu *CPU, d int, value uint8) {
	MemoryWrite(cpu, uint16(int(cpu.Reg.IX)+d), value)
}

func MEM_IY(cpu *CPU, d int) uint8 {
	return MemoryRead(cpu, uint16(int(cpu.Reg.IY)+d))
}

func MEM_IY_W(cpu *CPU, d int, value uint8) {
	MemoryWrite(cpu, uint16(int(cpu.Reg.IY)+d), value)
}

func MEM_BC(cpu *CPU) uint8 {
	return MemoryRead(cpu, cpu.Reg.BC())
}

func MEM_BC_W(cpu *CPU, value uint8) {
	MemoryWrite(cpu, cpu.Reg.BC(), value)
}

func MEM_DE(cpu *CPU) uint8 {
	return MemoryRead(cpu, cpu.Reg.DE())
}

func MEM_DE_W(cpu *CPU, value uint8) {
	MemoryWrite(cpu, cpu.Reg.DE(), value)
}

func MEM_SP(cpu *CPU, off uint8) uint8 {
	return MemoryRead(cpu, cpu.Reg.SP+uint16(off))
}

func MEM_SP_W(cpu *CPU, off uint8, value uint8) {
	MemoryWrite(cpu, cpu.Reg.SP+uint16(off), value)
}

func MemoryRead(cpu *CPU, addr uint16) uint8 {
	cpu.Pin.ADDR = addr
	cpu.Pin.RD = true
	cpu.Pin.MREQ = true

	cpu.Pin.Bus()

	if cpu.dataBreakPoints != nil {
		if cb, ok := cpu.dataBreakPoints[cpu.Pin.ADDR]; ok {
			cb(false, cpu.Pin.DATA)
		}
	}

	cpu.Pin.RD = false
	cpu.Pin.MREQ = false

	return cpu.Pin.DATA
}

func MemoryWrite(cpu *CPU, addr uint16, value uint8) {
	cpu.Pin.ADDR = addr
	cpu.Pin.DATA = value
	cpu.Pin.WR = true
	cpu.Pin.MREQ = true

	if cpu.dataBreakPoints != nil {
		if cb, ok := cpu.dataBreakPoints[cpu.Pin.ADDR]; ok {
			cb(true, cpu.Pin.DATA)
		}
	}

	cpu.Pin.Bus()

	cpu.Pin.WR = false
	cpu.Pin.MREQ = false
}

func MemoryRead16(cpu *CPU, addr uint16) uint16 {
	l := MemoryRead(cpu, addr)
	h := MemoryRead(cpu, addr+1)

	return helpers.To16(l, h)
}

func MemoryWrite16(cpu *CPU, addr uint16, value uint16) {
	l, h := helpers.To8(value)

	MemoryWrite(cpu, addr, l)
	MemoryWrite(cpu, addr+1, h)
}

func PushStack16(cpu *CPU, value uint16) {
	l, h := helpers.To8(value)

	cpu.Reg.SP--
	MemoryWrite(cpu, cpu.Reg.SP, h)
	cpu.Reg.SP--
	MemoryWrite(cpu, cpu.Reg.SP, l)
}

func PopStack16(cpu *CPU) uint16 {
	l := MemoryRead(cpu, cpu.Reg.SP)
	cpu.Reg.SP++
	h := MemoryRead(cpu, cpu.Reg.SP)
	cpu.Reg.SP++

	return helpers.To16(l, h)
}
