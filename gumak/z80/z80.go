package z80

import (
	"fmt"
	"mutex/gumak/log"
)

type Pins struct {
	// Bus
	ADDR uint16 // <->
	DATA uint8  // ->
	// System control
	M1    bool // ->
	MREQ  bool // ->
	IOREQ bool // ->
	RD    bool // ->
	WR    bool // ->
	RFSH  bool // ->
	// CPU bus control
	BUSREQ bool // <-
	BUSACK bool // ->
	// CPU control
	HALT  bool // ->
	WAIT  bool // <-
	INT   bool // <-
	NMI   bool // <-
	RESET bool // <-

	Bus func()
}

type CPU struct {
	Reg  Registers
	Pin  Pins
	IFF1 bool // Enable interrupts flip-flop1
	IFF2 bool // Enable interrupts flip-flop2

	maskableSkip  int  // After EI call maskable interrupts are disabled for next instruction (in case RETN)
	halted        bool // Halted after HALT call, in which case CPU is NOPing until interupt
	InterruptMode int

	// TODO: Remove?
	Frequency       int
	TStatesPerFrame int
	TStateUs        float64

	breakPoints     map[uint16]func()
	dataBreakPoints map[uint16]func(bool, uint8)

	addressCache []int
	symbols      *map[uint16]string
}

func NOP(cpu *CPU) int {
	log.Trace(2, "NOP")
	return 4
}

func HandleNMI(cpu *CPU) {
	log.Trace(1, "NMI handler")

	cpu.IFF1 = false
	PushStack16(cpu, cpu.Reg.PC)
	cpu.Reg.PC = 0x0066
}

func HandleInterrupt(cpu *CPU) int {
	log.Trace(1, "Interrupt handler [MOD %d]", cpu.InterruptMode)
	cpu.halted = false

	PushStack16(cpu, cpu.Reg.PC)
	cpu.IFF1 = false
	cpu.IFF2 = false

	switch cpu.InterruptMode {
	case 0:
		cpu.Reg.PC = 0x38
		return 13
	case 1:
		cpu.Reg.PC = 0x38
		return 13
	case 2:
		cpu.Reg.PC = MemoryRead16(cpu, uint16(cpu.Reg.I)<<8|0x00ff)
		return 19
	}

	panic(fmt.Sprintf("Invalid interrupt mode: %d", cpu.InterruptMode))
}

func (cpu *CPU) Tick() int {
	tStates := 0

	if cpu.Pin.NMI {
		HandleNMI(cpu)
	}

	if cpu.Pin.INT && (cpu.IFF1 && cpu.maskableSkip == 0) {
		tStates += HandleInterrupt(cpu)
	}

	if cpu.maskableSkip > 0 {
		cpu.maskableSkip--
	}

	cpu.Refresh(1)
	if cpu.halted {
		tStates += NOP(cpu)
	} else {
		tStates += DecodeAndExecute(cpu)
	}

	return tStates
}

func (cpu *CPU) AttachBreakpointAddr(addr uint16, cb func()) {
	if cpu.breakPoints == nil {
		cpu.breakPoints = make(map[uint16]func())
	}
	cpu.breakPoints[addr] = cb
}

func (cpu *CPU) AttachBreakpointName(name string, offset uint16, cb func()) {
	if cpu.dataBreakPoints == nil {
		cpu.dataBreakPoints = make(map[uint16]func(bool, uint8))
	}

	if cpu.symbols != nil {
		for a, n := range *cpu.symbols {
			if n == name {
				cpu.breakPoints[a+offset] = cb
				return
			}
		}
	}

	panic(fmt.Sprintf("Breakpoint '%s' not found", name))
}

func (cpu *CPU) DettachBreakpoint(addr uint16) {
	delete(cpu.breakPoints, addr)
	if len(cpu.breakPoints) == 0 {
		cpu.breakPoints = nil
	}
}

func (cpu *CPU) AttachBreakpointData(addr uint16, size uint16, cb func(bool, uint8)) {
	if cpu.dataBreakPoints == nil {
		cpu.dataBreakPoints = make(map[uint16]func(bool, uint8))
	}

	for i := uint16(0); i < size; i++ {
		cpu.dataBreakPoints[addr+i] = cb
	}
}

func (cpu *CPU) DettachBreakpointData(addr uint16, size uint16) {
	for i := uint16(0); i < size; i++ {
		delete(cpu.dataBreakPoints, addr+i)
	}

	if len(cpu.dataBreakPoints) == 0 {
		cpu.dataBreakPoints = nil
	}
}

func (cpu *CPU) Restart() {
	cpu.Init(cpu.Frequency, cpu.TStatesPerFrame, cpu.symbols)
}

func (cpu *CPU) Reset() {
	cpu.IFF1 = false
	cpu.IFF2 = false
	cpu.halted = false
	cpu.maskableSkip = 0
}

func (cpu *CPU) Refresh(v uint8) {
	cpu.Reg.R += v
}

func (cpu *CPU) Init(frequency int, tStatePerFrame int, symbols *map[uint16]string) {
	log.Info("=== Z80 initialize ===")

	cpu.Pin.ADDR = 0
	cpu.Pin.DATA = 0

	cpu.Pin.M1 = false
	cpu.Pin.MREQ = false
	cpu.Pin.IOREQ = false
	cpu.Pin.RD = false
	cpu.Pin.WR = false
	cpu.Pin.RFSH = false
	cpu.Pin.BUSREQ = false
	cpu.Pin.BUSACK = false
	cpu.Pin.HALT = false
	cpu.Pin.WAIT = false
	cpu.Pin.INT = false
	cpu.Pin.NMI = false
	cpu.Pin.RESET = false

	cpu.Frequency = frequency
	cpu.TStatesPerFrame = tStatePerFrame
	cpu.TStateUs = 1e6 / float64(cpu.Frequency)

	cpu.Reg.Clear()

	cpu.Reset()

	cpu.symbols = symbols
}
