package device

import (
	"embed"
	"mutex/gumak/log"
)

type Ram struct {
	roms [2][0x4000]uint8
	// 8 banks of 16K.
	banks [8][0x4000]uint8

	// Mapping of active banks.
	page [4][]uint8

	pagingEnabled bool

	activeRom  int
	activeBank int
}

const (
	BANK_VRAM        = 5
	BANK_VRAM_SHADOW = 7
)

func (r *Ram) DumpState() {
	log.Debug("ROM / PAGE[0] = %d ", r.activeRom)
	log.Debug("PAGE[3] = BANK_%d ", r.activeBank)
}

func (r *Ram) SetRom(rom int) {
	if r.pagingEnabled {
		r.activeRom = rom
		r.page[0] = r.roms[rom][:]
	}
}

func (r *Ram) SetPageBank(page int, bank int) {
	if r.pagingEnabled {
		if page != 3 {
			panic("Swapping other bank than 3!")
		}
		r.activeBank = bank
		r.page[page] = r.banks[bank][:]
	}
}

func (r *Ram) Init() {
	r.pagingEnabled = true

	r.SetRom(0)
	r.page[1] = r.banks[5][:]
	r.page[2] = r.banks[2][:]
	r.page[3] = r.banks[0][:]
}

func (r *Ram) SetPagingEnabled(enabled bool) {
	r.pagingEnabled = enabled
}

func (r *Ram) Page(page int) []uint8 {
	return r.page[page]
}

func (r *Ram) Bank(bank int) []uint8 {
	return r.banks[bank][:]
}

func (r *Ram) Rom(rom int) []uint8 {
	return r.roms[rom][:]
}

func (r *Ram) SetRomContent(rom int, bytes []byte, offset uint16) {
	if len(bytes) > 0x4000-int(offset) {
		panic("Invalid page size")
	}

	for i, b := range bytes {
		r.roms[rom][int(offset)+i] = b
	}
}

func (r *Ram) SetPageContent(page int, bytes []byte, offset uint16) {
	if len(bytes) > 0x4000-int(offset) {
		panic("Invalid page size")
	}

	for i, b := range bytes {
		r.page[page][int(offset)+i] = b
	}
}

func (r *Ram) SetBankContent(bank int, bytes []byte, offset uint16) {
	if len(bytes) > 0x4000-int(offset) {
		panic("Invalid bank size")
	}

	for i, b := range bytes {
		r.banks[bank][int(offset)+i] = b
	}
}

func (r *Ram) Size() uint32 {
	return uint32(64 * 1024)
}

func (r *Ram) LoadRom(fs embed.FS, file string, rom int, length uint16) (int, error) {
	f, err := embed.FS.Open(fs, file)
	if err != nil {
		return 0, err
	}

	size, err := f.Read(r.roms[rom][:length])
	if err != nil {
		return 0, err
	}

	log.Info("Loaded %d bytes into ROM[%d]", size, rom)

	return size, nil
}

func pageOffset(addr uint16) (page int, offset uint16) {
	return int(addr) >> 14, addr & 0x3fff
}

func (r *Ram) Read(addr uint16) uint8 {
	page, offset := pageOffset(addr)
	return r.page[page][offset]
}

func (r *Ram) Write(addr uint16, dataBus uint8) {
	page, offset := pageOffset(addr)
	r.page[page][offset] = dataBus
}
