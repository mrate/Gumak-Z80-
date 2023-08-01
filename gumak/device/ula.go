package device

import "mutex/gumak/helpers"

// FE port:
//
//    +---+------+---+-------+------+---+---+---+
//    | 7 |  6   | 5 |   4   |  3   | 2 | 1 | 0 |
//    +---+------+---+-------+------+---+---+---+
// RD | X | tape | X |         keyboard         |
//    +---+------+---+-------+------+-----------+
// WR | X |  X   | X | sound | tape |   border  |
//    +---+------+---+-------+------+-----------+

//  IN:    Reads keys (bit 0 to bit 4 inclusive)
//
// 0xfefe  SHIFT, Z, X, C, V            0xeffe  0, 9, 8, 7, 6
// 0xfdfe  A, S, D, F, G                0xdffe  P, O, I, U, Y
// 0xfbfe  Q, W, E, R, T                0xbffe  ENTER, L, K, J, H
// 0xf7fe  1, 2, 3, 4, 5                0x7ffe  SPACE, SYM SHFT, M, N, B

type Ula struct {
	Beeper    *Beeper
	Ay_3_8912 *AY_3_8912
	Tape      *Tape

	ram *Ram

	Frames int

	Is128K     bool
	Last0x7ffd uint8

	VRamBank    int
	BorderColor uint8

	Keyboard [8]uint8
}

func (ula *Ula) Init(ram *Ram, beeper *Beeper, sound *AY_3_8912, is128k bool) {
	ula.Beeper = beeper
	ula.Ay_3_8912 = sound
	ula.Tape = &Tape{}
	ula.ram = ram
	ula.Is128K = is128k

	ula.VRamBank = BANK_VRAM

	for i := range ula.Keyboard {
		ula.Keyboard[i] = 0b11111
	}
}

func (ula *Ula) Write(addr uint16, value uint8) {
	al, ah := helpers.To8(addr)

	switch al {
	case 0xfe:
		ula.BorderColor = value & 0b00000111
		if !ula.Tape.Running {
			bit := (value & 0b10000) != 0
			ula.Beeper.Beep(bit)
			ula.Tape.earBit = bit
		}
	case 0xfd:
		if ula.Is128K {
			switch ah {
			case 0x7f:
				ula.Last0x7ffd = value

				// Bit 0-2: bank select for page 0xc000
				ula.ram.SetPageBank(3, int(value&0b111))
				// Bit 3: 0=VRAM, 1=Shadow VRAM
				if value&0b1000 == 0 {
					ula.VRamBank = BANK_VRAM
				} else {
					ula.VRamBank = BANK_VRAM_SHADOW
				}
				// Bit 4: ROM select
				ula.ram.SetRom(int(value&0b10000) >> 4)
				// Bit 5: Disable pagging until reset.
				if value&0b100000 != 0 {
					ula.ram.SetPagingEnabled(false)
				}

			case 0xff: // AY-3-8912 sound chip
				ula.Ay_3_8912.SelectRegister(value)
			case 0xbf: // AY-3-8912 sound chip
				ula.Ay_3_8912.Write(value)
			}
		}
	}
}

func (ula *Ula) Read(addr uint16) uint8 {
	_, ah := helpers.To8(addr)

	b := uint8(0)

	switch ah {
	case 0xff:
		b = ula.Ay_3_8912.Read()
	case 0xfe:
		b |= ula.Keyboard[0]
	case 0xfd:
		b |= ula.Keyboard[1]
	case 0xfb:
		b |= ula.Keyboard[2]
	case 0xf7:
		b |= ula.Keyboard[3]
	case 0xef:
		b |= ula.Keyboard[4]
	case 0xdf:
		b |= ula.Keyboard[5]
	case 0xbf:
		b |= ula.Keyboard[6]
	case 0x7f:
		b |= 0b10100000
		if ula.Tape.EarBit() {
			b |= (1 << 6)
		}
		b |= ula.Keyboard[7]
	default:
		b = 0xff
	}

	return b
}

func (ula *Ula) UpdateEndFrame() {
	ula.Frames = (ula.Frames + 1) % 32
}
