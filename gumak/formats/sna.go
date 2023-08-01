package formats

import (
	"errors"
	"io"
	"mutex/gumak/device"
	"mutex/gumak/helpers"
	"mutex/gumak/z80"
)

type SNA struct {
	header [27]byte
	ram    [48 * 1024]byte
}

func (s *SNA) Load(reader io.Reader, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error {
	size, err := reader.Read(s.header[:])
	if err != nil {
		return nil
	}

	if size < len(s.header) {
		return errors.New("Invalid header")
	}

	size, err = reader.Read(s.ram[:])
	if err != nil {
		return nil
	}

	if size < len(s.ram) {
		return errors.New("Invalid RAM")
	}

	cpu.Reset()

	cpu.Reg.I = s.header[0]
	cpu.Reg.L_ = s.header[1]
	cpu.Reg.H_ = s.header[2]
	cpu.Reg.E_ = s.header[3]
	cpu.Reg.D_ = s.header[4]
	cpu.Reg.C_ = s.header[5]
	cpu.Reg.B_ = s.header[6]
	cpu.Reg.F_ = s.header[7]
	cpu.Reg.A_ = s.header[8]
	cpu.Reg.L = s.header[9]
	cpu.Reg.H = s.header[10]
	cpu.Reg.E = s.header[11]
	cpu.Reg.D = s.header[12]
	cpu.Reg.C = s.header[13]
	cpu.Reg.B = s.header[14]
	cpu.Reg.IY = helpers.To16(s.header[15], s.header[16])
	cpu.Reg.IX = helpers.To16(s.header[17], s.header[18])
	if s.header[19]&0b100 != 0 {
		cpu.IFF2 = true
	} else {
		cpu.IFF2 = false
	}
	cpu.IFF1 = cpu.IFF2
	cpu.Reg.R = s.header[20]
	cpu.Reg.F = s.header[21]
	cpu.Reg.A = s.header[22]
	cpu.Reg.SP = helpers.To16(s.header[23], s.header[24])
	cpu.InterruptMode = int(s.header[25])
	ula.BorderColor = s.header[26] & 0b111

	ram.SetPageContent(1, s.ram[0:0x4000], 0)
	ram.SetPageContent(2, s.ram[0x4000:0x8000], 0)
	ram.SetPageContent(3, s.ram[0x8000:], 0)

	cpu.Reg.PC = 0x72

	return nil
}

func (s *SNA) Save(writer io.Writer, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error {
	s.header[0] = cpu.Reg.I
	s.header[1] = cpu.Reg.L_
	s.header[2] = cpu.Reg.H_
	s.header[3] = cpu.Reg.E_
	s.header[4] = cpu.Reg.D_
	s.header[5] = cpu.Reg.C_
	s.header[6] = cpu.Reg.B_
	s.header[7] = cpu.Reg.F_
	s.header[8] = cpu.Reg.A_
	s.header[9] = cpu.Reg.L
	s.header[10] = cpu.Reg.H
	s.header[11] = cpu.Reg.E
	s.header[12] = cpu.Reg.D
	s.header[13] = cpu.Reg.C
	s.header[14] = cpu.Reg.B
	s.header[15], s.header[16] = helpers.To8(cpu.Reg.IY)
	s.header[17], s.header[18] = helpers.To8(cpu.Reg.IX)
	s.header[19] = 0
	if cpu.IFF2 {
		s.header[19] = 0b100
	}
	s.header[20] = cpu.Reg.R
	s.header[21] = cpu.Reg.F
	s.header[22] = cpu.Reg.A
	s.header[23], s.header[24] = helpers.To8(cpu.Reg.SP)
	s.header[25] = uint8(cpu.InterruptMode)
	s.header[26] = ula.BorderColor

	// TODO: Error
	size, _ := writer.Write(s.header[:])
	if size < len(s.header) {
		return errors.New("Error writing header")
	}

	for page := 1; page < 4; page++ {
		ramData := ram.Page(page)
		// TODO: Error
		size, _ = writer.Write(ramData)
		if size < len(ramData) {
			return errors.New("Error writing ram")
		}
	}

	return nil
}
