package z80

func ReadIO(cpu *CPU, addr uint16) uint8 {
	cpu.Pin.ADDR = addr
	cpu.Pin.RD = true
	cpu.Pin.IOREQ = true

	cpu.Pin.Bus()

	cpu.Pin.RD = false
	cpu.Pin.IOREQ = false

	return cpu.Pin.DATA
}

func WriteIO(cpu *CPU, addr uint16, value uint8) {
	cpu.Pin.ADDR = addr
	cpu.Pin.DATA = value
	cpu.Pin.RD = false
	cpu.Pin.IOREQ = true

	cpu.Pin.Bus()

	cpu.Pin.IOREQ = false
}
