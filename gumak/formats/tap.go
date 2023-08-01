package formats

import (
	"bytes"
	"errors"
	"io"
	"mutex/gumak/helpers"
)

type Block interface {
	Length() int
	Data() []byte
}

type blockHeader struct {
	data           []byte
	typ            byte
	name           string
	length         uint16
	param1, param2 uint16
}

type blockData []byte

type TAP struct {
	filename string
	blocks   []Block
	data     []byte
}

func (h *blockHeader) Length() int  { return len(h.data) }
func (h *blockHeader) Data() []byte { return h.data }
func (h blockData) Length() int     { return len(h) }
func (h blockData) Data() []byte    { return h }

func (t *TAP) Byte(position int) uint8    { return t.data[position] }
func (t *TAP) Length() int                { return len(t.data) }
func (t *TAP) BlockLength(block int) int  { return t.blocks[block].Length() }
func (t *TAP) BlockData(block int) []byte { return t.blocks[block].Data() }
func (t *TAP) BlockCount() int            { return len(t.blocks) }
func (t *TAP) IsHeader(block int) bool    { return t.blocks[block].Data()[0] == 0x0 }

func (tap *TAP) readInt16(reader io.Reader) (uint16, error) {
	var length []byte = make([]byte, 2)
	size, err := reader.Read(length)

	if size < 2 || err != nil {
		return 0, err
	}

	value := helpers.To16(length[0], length[1])
	return value, nil
}

func (tap *TAP) readHeader(data []byte) Block {
	header := &blockHeader{}

	header.data = data
	header.typ = data[1]
	header.name = string(data[2:12])
	header.length = helpers.To16(data[12], data[13])
	header.param1 = helpers.To16(data[14], data[15])
	header.param2 = helpers.To16(data[16], data[17])

	return header
}

func calcChecksum(data []byte) bool {
	sum := uint8(0)
	for _, v := range data {
		sum ^= v
	}
	return sum == 0
}

func (tap *TAP) Read(p []byte) error {
	reader := bytes.NewReader(p)

	dataLength := 0
	for {
		length, err := tap.readInt16(reader)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		if err != nil {
			return err
		}

		block := make([]byte, length)
		len, err := reader.Read(block)
		if err != nil {
			return err
		}

		if len < int(length) {
			return errors.New("Unexpected EOF")
		}

		if !calcChecksum(block) {
			return errors.New("Invalid checksum")
		}

		if block[0] == 0x0 {
			// Header.
			tap.blocks = append(tap.blocks, tap.readHeader(block))
		} else {
			// Data.
			tap.blocks = append(tap.blocks, blockData(block))
		}

		dataLength += int(length)
	}

	tap.data = make([]byte, dataLength)
	c := 0
	for _, block := range tap.blocks {
		for _, data := range block.Data() {
			tap.data[c] = data
			c++
		}
	}

	return nil
}
