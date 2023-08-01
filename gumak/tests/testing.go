package tests

import (
	"mutex/gumak/device"
	"mutex/gumak/z80"
)

type Hardware struct {
	cpu z80.CPU
	ram device.Ram
}

func TestHw() *Hardware {
	hw := new(Hardware)
	hw.cpu.Init(35000000, 224*312, nil)
	hw.cpu.Reset()
	hw.ram.Init()

	hw.cpu.Pin.Bus = func() {
		if hw.cpu.Pin.MREQ {
			// Memory request
			if hw.cpu.Pin.RD {
				hw.cpu.Pin.DATA = hw.ram.Read(hw.cpu.Pin.ADDR)
			} else if hw.cpu.Pin.WR {
				hw.ram.Write(hw.cpu.Pin.ADDR, hw.cpu.Pin.DATA)
			}
		}
	}

	return hw
}
