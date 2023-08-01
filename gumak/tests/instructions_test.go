package tests

import (
	"mutex/gumak/z80"
	"testing"
)

func TestInstrCount(t *testing.T) {
	l := len(z80.OpCodes)
	if l != 256 {
		t.Fatalf("Invalid OpCodes len: %d", len(z80.OpCodes))
	}
}

func TestNeg(t *testing.T) {
	hw := TestHw()

	hw.cpu.Reg.A = 0b10011000

	hw.ram.Write(0x0000, 0xed)
	hw.ram.Write(0x0001, 0x44)

	instr := z80.DecodeInstruction(&hw.cpu)

	instr(&hw.cpu)

	if hw.cpu.Reg.A != 0b01101000 {
		t.Fatalf("Expected A=%x, got A=%x", 0b01101000, hw.cpu.Reg.A)
	}
}

func TestADD_HL_ss(t *testing.T) {
	hw := TestHw()

	hw.cpu.Reg.HL_write(0x4242)
	hw.cpu.Reg.DE_write(0x1111)
	hw.cpu.Reg.PC = 0

	hw.ram.Write(0x0000, 0b00011001)
	instr := z80.DecodeInstruction(&hw.cpu)

	instr(&hw.cpu)

	if hw.cpu.Reg.HL() != 0x5353 {
		t.Fatalf("Expected HL=%x, got HL=%x", 0x5353, hw.cpu.Reg.HL())
	}
}

func TestSET_B_IX(t *testing.T) {
	hw := TestHw()

	hw.cpu.Reg.IX = 0x2000
	hw.ram.Write(0x2003, 0b10101010)

	hw.ram.Write(0x0000, 0xdd)
	hw.ram.Write(0x0001, 0xcb)
	hw.ram.Write(0x0002, 0x03)
	hw.ram.Write(0x0003, 0b11000110)

	hw.cpu.Reg.PC = 0

	instr := z80.DecodeInstruction(&hw.cpu)

	instr(&hw.cpu)

	data := hw.ram.Read(0x2003)

	if data != 0b10101011 {
		t.Fatalf("Expected HL=%x, got HL=%x", 0b10101011, data)
	}

}
