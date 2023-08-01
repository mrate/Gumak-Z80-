package formats

import (
	"errors"
	"io"
	"mutex/gumak/device"
	"mutex/gumak/z80"
	"path/filepath"
)

type SnapshotLoader interface {
	Load(reader io.Reader, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error
	Save(writer io.Writer, cpu *z80.CPU, ula *device.Ula, ram *device.Ram) error
}

func NewSnapshot(filename string) (SnapshotLoader, error) {
	switch filepath.Ext(filename) {
	case ".sna":
		return &SNA{}, nil
	case ".z80":
		return &Z80{}, nil
	}

	return nil, errors.New("Invalid snapshot format")
}

func NewTape(filename string) (device.TapeSource, error) {
	switch filepath.Ext(filename) {
	case ".tap":
		return &TAP{filename: filename}, nil
	}

	return nil, errors.New("Invalid tape format")
}
