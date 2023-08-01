package formats

import (
	"errors"
	"fmt"
	"io"
	"mutex/gumak/device"
	"mutex/gumak/helpers"
	"mutex/gumak/z80"
)

type Z80 struct {
	version    int
	compressed bool
	header     [30]byte
	header2    []byte
	is128k     bool
}

func decompressedSize(compressed []byte) int {
	size, ed := 0, 0
	for _, b := range compressed {
		switch ed {
		case 2:
			size += int(b)
			ed++
		case 3:
			ed = 0
		default:
			if b == 0xed {
				ed++
			} else if ed == 1 {
				size += 2
				ed = 0
			} else {
				size++
			}
		}
	}

	return size
}

func putByte(output []byte, value byte, pos *int) {
	output[*pos] = value
	*pos++
}

func decompress(compressed []byte) ([]byte, error) {
	size := decompressedSize(compressed)
	decompressed := make([]byte, size)

	ed, out, count := 0, 0, 0

	for _, b := range compressed {
		switch ed {
		case 0:
			if b == 0xed {
				ed++
			} else {
				putByte(decompressed, b, &out)
			}
		case 1:
			if b == 0xed {
				ed++
			} else {
				ed = 0
				putByte(decompressed, 0xed, &out)
				putByte(decompressed, b, &out)
			}
		case 2:
			ed++
			count = int(b)
		case 3:
			ed = 0
			for i := 0; i < count; i++ {
				putByte(decompressed, b, &out)
			}
		}
	}

	if out != size {
		return nil, errors.New("Invalid compressed data")
	}

	return decompressed, nil
}

func compress(raw []byte) []byte {
	panic("Implement me!")
}

func (s *Z80) loadVersion1(f io.Reader, ram *device.Ram) error {
	data := make([]byte, 0xc000) // 48K
	size, err := f.Read(data)
	if err != nil {
		return err
	}

	data = data[:size]

	end := data[size-4:]
	if end[0] != 0 || end[1] != 0xed || end[2] != 0xed || end[3] != 0 {
		return errors.New("Invalid end block")
	}

	if s.compressed {
		data, err = decompress(data[:size-4])
	}

	if len(data) != 0xc000 {
		return errors.New("Invalid snapshot ram size")
	}

	ram.SetPageContent(1, data[0:0x4000], 0)
	ram.SetPageContent(2, data[0x4000:0x8000], 0)
	ram.SetPageContent(3, data[0x8000:0xc000], 0)

	return nil
}

func (s *Z80) loadVersion23(f io.Reader, ram *device.Ram) error {
	// Memory blocks:
	// Byte    Length  Description
	//    ---------------------------
	//    0       2       Length of compressed data (without this 3-byte header)
	//                    If length=0xffff, data is 16384 bytes long and not compressed
	//    2       1       Page number of block
	//    3       [0]     Data
	for {
		var blockHeader [3]byte
		size, err := f.Read(blockHeader[:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if size == 0 {
			return nil
		}

		if size < len(blockHeader) {
			return errors.New("Unexpected end of file")
		}

		compressed := s.compressed
		length := int(helpers.To16(blockHeader[0], blockHeader[1]))
		if length == 0xffff {
			length = 16 * 1024
			compressed = false
		}

		page := blockHeader[2]

		compressedData := make([]byte, length)
		size, err = f.Read(compressedData)
		if err != nil {
			return err
		}

		if size != length {
			return errors.New("Unexpected end of file")
		}

		data := compressedData
		if compressed {
			data, err = decompress(compressedData)
			if err != nil {
				return err
			}
		}

		//  The pages are numbered, depending on the hardware mode, in the following way:
		//
		//     Page    In '48 mode     In '128 mode    In SamRam mode
		//     ------------------------------------------------------
		//      0      48K rom         rom (basic)     48K rom
		//      1      Interface I, Disciple or Plus D rom, according to setting
		//      2      -               rom (reset)     samram rom (basic)
		//      3      -               page 0          samram rom (monitor,..)
		//      4      8000-bfff       page 1          Normal 8000-bfff
		//      5      c000-ffff       page 2          Normal c000-ffff
		//      6      -               page 3          Shadow 8000-bfff
		//      7      -               page 4          Shadow c000-ffff
		//      8      4000-7fff       page 5          4000-7fff
		//      9      -               page 6          -
		//     10      -               page 7          -
		//     11      Multiface rom   Multiface rom   -
		// In 48K mode, pages 4,5 and 8 are saved. In SamRam mode, pages 4 to 8 are saved. In 128K mode, all pages from 3 to 10 are saved. Pentagon snapshots are very similar to 128K snapshots, while Scorpion snapshots have the 16 RAM pages saved in pages 3 to 18. There is no end marker.
		if s.is128k {
			// 128K
			switch page {
			case 0, 2:
				ram.SetRomContent(int(page>>1), data, 0)
			case 3, 4, 5, 6, 7, 8, 9, 10:
				ram.SetBankContent(int(page-3), data, 0)
			}
		} else {
			// 48K
			switch page {
			case 0:
				ram.SetRomContent(0, data, 0)
			case 4:
				ram.SetPageContent(2, data, 0)
			case 5:
				ram.SetPageContent(3, data, 0)
			case 8:
				ram.SetPageContent(1, data, 0)
			}
		}
	}
}

func (s *Z80) saveBlock(f io.Writer, page int, block []byte, compression bool) error {
	var header [3]byte
	if compression {
		header[0], header[1] = helpers.To8(uint16(len(block)))
	} else {
		header[0], header[1] = 0xff, 0xff
	}

	header[2] = uint8(page)

	size, err := f.Write(header[:])
	if err != nil {
		return err
	}

	if size < len(header) {
		return errors.New("Error writing snapshot")
	}

	data := block
	if compression {
		data = compress(block)
	}

	size, err = f.Write(data)
	if err != nil {
		return err
	}

	if size != len(data) {
		return errors.New("Error writing snapshot")
	}

	return nil
}

func (s *Z80) saveVersion2(f io.Writer, ram *device.Ram) error {
	comprim := false

	if s.is128k {
		s.saveBlock(f, 0, ram.Rom(0), comprim)
		s.saveBlock(f, 2, ram.Rom(1), comprim)

		for i := 0; i < 8; i++ {
			s.saveBlock(f, i+3, ram.Bank(i), comprim)
		}
	} else {
		s.saveBlock(f, 0, ram.Rom(0), comprim)
		s.saveBlock(f, 4, ram.Page(2), comprim)
		s.saveBlock(f, 5, ram.Page(3), comprim)
		s.saveBlock(f, 8, ram.Page(1), comprim)
	}

	return nil
}

func (s *Z80) Load(reader io.Reader, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error {
	size, err := reader.Read(s.header[:])
	if err != nil {
		return nil
	}

	if size < len(s.header) {
		return errors.New("Invalid header")
	}

	//        0       1       A register
	//        1       1       F register
	//        2       2       BC register pair (LSB, i.e. C, first)
	//        4       2       HL register pair
	//        6       2       Program counter
	//        8       2       Stack pointer
	//        10      1       Interrupt register
	//        11      1       Refresh register (Bit 7 is not significant!)
	//        12      1       Bit 0  : Bit 7 of the R-register
	//                        Bit 1-3: Border colour
	//                        Bit 4  : 1=Basic SamRom switched in
	//                        Bit 5  : 1=Block of data is compressed
	//                        Bit 6-7: No meaning
	//        13      2       DE register pair
	//        15      2       BC' register pair
	//        17      2       DE' register pair
	//        19      2       HL' register pair
	//        21      1       A' register
	//        22      1       F' register
	//        23      2       IY register (Again LSB first)
	//        25      2       IX register
	//        27      1       Interrupt flipflop, 0=DI, otherwise EI
	//        28      1       IFF2 (not particularly important...)
	//        29      1       Bit 0-1: Interrupt mode (0, 1 or 2)
	//                        Bit 2  : 1=Issue 2 emulation
	//                        Bit 3  : 1=Double interrupt frequency
	//                        Bit 4-5: 1=High video synchronisation
	//                                 3=Low video synchronisation
	//                                 0,2=Normal
	//                        Bit 6-7: 0=Cursor/Protek/AGF joystick
	//                                 1=Kempston joystick
	//                                 2=Sinclair 2 Left joystick (or user
	//                                   defined, for version 3 .z80 files)
	//                                 3=Sinclair 2 Right joystick
	//

	ram.Init()

	s.version = 1
	s.is128k = false

	cpu.Reg.A = s.header[0]
	cpu.Reg.F = s.header[1]
	cpu.Reg.B, cpu.Reg.C = s.header[3], s.header[2]
	cpu.Reg.H, cpu.Reg.L = s.header[5], s.header[4]
	cpu.Reg.PC = helpers.To16(s.header[6], s.header[7])
	cpu.Reg.SP = helpers.To16(s.header[8], s.header[9])
	cpu.Reg.I = s.header[10]
	cpu.Reg.R = s.header[11] & 0b1111111

	cpu.Reg.R7 = s.header[12] & 0b1
	borderColor := s.header[12] & 0b1110 >> 1
	//samRom := (s.header[12] & 0b10000) != 0

	// TODO:
	//s.compressed = (s.header[12] & 0b100000) != 0
	s.compressed = true

	ula.BorderColor = borderColor

	cpu.Reg.D, cpu.Reg.E = s.header[14], s.header[13]
	cpu.Reg.B_, cpu.Reg.C_ = s.header[16], s.header[15]
	cpu.Reg.D_, cpu.Reg.E_ = s.header[18], s.header[17]
	cpu.Reg.H_, cpu.Reg.L_ = s.header[20], s.header[19]
	cpu.Reg.A_ = s.header[21]
	cpu.Reg.F_ = s.header[22]
	cpu.Reg.IY = helpers.To16(s.header[23], s.header[24])
	cpu.Reg.IX = helpers.To16(s.header[25], s.header[26])
	cpu.IFF1 = s.header[27] > 0
	cpu.IFF2 = s.header[28] > 0

	cpu.InterruptMode = int(s.header[29] & 0b11)
	//issue2 = int(s.header[29]&0b100) != 0
	//doubleIntFreq = int(s.header[29]&0b1000) != 0
	//videoSync = int(s.header[29]&0b110000) >> 4
	//joystick = int(s.header[29]&0b11000000) >> 6

	if cpu.Reg.PC == 0 {
		// Version 2 or 3

		// Offset  Length  Description
		//
		// * 30      2       Length of additional header block (see below)
		// * 32      2       Program counter
		// * 34      1       Hardware mode (see below)
		// * 35      1       If in SamRam mode, bitwise state of 74ls259.
		// 				  For example, bit 6=1 after an OUT 31,13 (=2*6+1)
		// 				  If in 128 mode, contains last OUT to 7ffd
		// * 36      1       Contains 0FF if Interface I rom paged
		// * 37      1       Bit 0: 1 if R register emulation on
		// 				  Bit 1: 1 if LDIR emulation on
		// * 38      1       Last OUT to fffd (soundchip register number)
		// * 39      16      Contents of the sound chip registers
		//   55      2       Low T state counter
		//   57      1       Hi T state counter
		//   58      1       Flag byte used by Spectator (QL spec. emulator)
		// 				  Ignored by Z80 when loading, zero when saving
		//   59      1       0FF if MGT Rom paged
		//   60      1       0FF if Multiface Rom paged. Should always be 0.
		//   61      1       0FF if 0-8191 is RAM
		//   62      1       0FF if 8192-16383 is RAM
		//   63      10      5x keyboard mappings for user defined joystick
		//   73      10      5x ascii word: keys corresponding to mappings above
		//   83      1       MGT type: 0=Disciple+Epson,1=Discipls+HP,16=Plus D
		//   84      1       Disciple inhibit button status: 0=out, 0ff=in
		//   85      1       Disciple inhibit flag: 0=rom pageable, 0ff=not

		var header2len [2]byte

		size, err = reader.Read(header2len[:])
		if err != nil {
			return err
		}

		if size != len(header2len) {
			return errors.New("Unexpected end of file")
		}

		// The value of the word at position 30 is 23 for version 2.01 files, and 54 for version 3.0 files
		len := int(helpers.To16(header2len[0], header2len[1]))
		switch len {
		case 23:
			s.version = 2
		case 54:
			s.version = 3
		case 55:
			s.version = 3
		}

		s.header2 = make([]byte, len)
		size, err = reader.Read(s.header2)
		if err != nil {
			return err
		}

		if size < len {
			return errors.New("Unexpected end of file")
		}

		cpu.Reg.PC = helpers.To16(s.header2[0], s.header2[1])

		// Value:          Meaning in v2.01        Meaning in v3.0x
		// --------------------------------------------------------
		// 0               48k                     48k
		// 1               48k + If.1              48k + If.1
		// 2               SamRam                  SamRam
		// 3               128k                    48k + M.G.T.
		// 4               128k + If.1             128k
		// 5               -                       128k + If.1
		// 6               -                       128k + M.G.T.
		if s.version == 2 {
			switch s.header2[2] {
			case 0:
				s.is128k = false
			case 3:
				s.is128k = true
			default:
				return errors.New(fmt.Sprintf("Unsupported hardware version: %d (file ver: %d)", s.header2[4], s.version))
			}
		} else if s.version == 3 {
			switch s.header2[2] {
			case 0:
				s.is128k = false
			case 4:
				s.is128k = true
			default:
				return errors.New(fmt.Sprintf("Unsupported hardware version: %d (file ver: %d)", s.header2[4], s.version))
			}
		}

		if s.is128k {
			out7ffd := s.header2[3]
			ula.Write(0x7ffd, out7ffd)
		}
	}

	if s.version == 1 {
		return s.loadVersion1(reader, ram)
	} else {
		return s.loadVersion23(reader, ram)
	}
}

func (s *Z80) Save(writer io.Writer, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error {
	s.compressed = false
	s.version = 2

	s.header[0] = cpu.Reg.A
	s.header[1] = cpu.Reg.F
	s.header[3], s.header[2] = cpu.Reg.B, cpu.Reg.C
	s.header[5], s.header[4] = cpu.Reg.H, cpu.Reg.L
	s.header[6], s.header[7] = 0, 0 // PC=0 - version 2, 3
	s.header[8], s.header[9] = helpers.To8(cpu.Reg.SP)
	s.header[10] = cpu.Reg.I
	s.header[11] = cpu.Reg.R & 0b1111111

	s.header[12] = cpu.Reg.R7 & 0b1

	s.header[12] = 0
	s.header[12] |= (ula.BorderColor & 0b111) << 1
	if s.compressed {
		s.header[12] |= 0b100000
	}

	s.header[14], s.header[13] = cpu.Reg.D, cpu.Reg.E
	s.header[16], s.header[15] = cpu.Reg.B_, cpu.Reg.C_
	s.header[18], s.header[17] = cpu.Reg.D_, cpu.Reg.E_
	s.header[20], s.header[19] = cpu.Reg.H_, cpu.Reg.L_
	s.header[21] = cpu.Reg.A_
	s.header[22] = cpu.Reg.F_
	s.header[23], s.header[24] = helpers.To8(cpu.Reg.IY)
	s.header[25], s.header[26] = helpers.To8(cpu.Reg.IX)
	if cpu.IFF1 {
		s.header[27] = 1
	}
	if cpu.IFF2 {
		s.header[28] = 1
	}

	s.header[29] = uint8(cpu.InterruptMode)

	s.header2 = make([]byte, 2+23)
	s.header2[0], s.header2[1] = helpers.To8(23)
	s.header2[2], s.header2[3] = helpers.To8(cpu.Reg.PC)

	if ula.Is128K {
		s.header2[4] = 3
		s.header[5] = ula.Last0x7ffd
	} else {
		s.header2[4] = 0
	}

	// TODO: Handle error
	size, _ := writer.Write(s.header[:])
	if size < len(s.header) {
		return errors.New("Error writing header")
	}

	size, _ = writer.Write(s.header2)
	if size < len(s.header2) {
		return errors.New("Error writing header")
	}

	return s.saveVersion2(writer, ram)
}
