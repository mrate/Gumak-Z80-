package device

// https://uelectronics.info/2015/03/21/zx-spectrum-and-loaders-part-one/

import (
	"errors"
	"mutex/gumak/log"
	"os"
	"time"
)

const (
	LEADER_PULSE_WIDTH = 2168

	TAPE_SYNC_HIGH_PULSE_WIDTH = 667
	TAPE_SYNC_LOW_PULSE_WIDTH  = 735

	TAPE_BIT_HIGH_PULSE_WIDTH = 1710
	TAPE_BIT_LOW_PULSE_WIDTH  = 855

	TAPE_HEADER_LEADER_PULSES = 8063
	TAPE_DATA_LEADER_PULSES   = 3223

	TAPE_PAUSE_PULSE_WIDTH = 3500000
)

const (
	TAPE_STATE_BLOCK_START = iota
	TYPE_STATE_BLOCK_LEADER
	TAPE_STATE_SYNC
	TAPE_STATE_BIT_HIGH
	TAPE_STATE_BIT_LOW
	TAPE_STATE_PAUSE
	TAPE_STATE_PAUSE_END
	TAPE_STATE_FINISHED
)

type TapeSource interface {
	Read(b []byte) error
	Byte(position int) uint8
	Length() int
	BlockData(block int) []byte
	BlockLength(block int) int
	BlockCount() int
	IsHeader(block int) bool
}

type Tape struct {
	data   TapeSource
	earBit bool
	state  int

	pulseCounter      int
	pulseWidthCounter int

	block       int
	position    int
	blockLength int
	bitMask     uint8

	Running bool

	start time.Time
}

func (tape *Tape) Load(file string, source TapeSource) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return errors.New("Failed to open file")
	}

	tape.data = source
	err = tape.data.Read(data)

	if err != nil {
		return err
	}

	log.Info("Loaded TAP: %s", file)
	return nil
}

func (tape *Tape) Start() {
	tape.Running = true
	tape.state = TAPE_STATE_BLOCK_START
	tape.position = 0
	tape.block = 0
	tape.start = time.Now()

	if tape.data.BlockCount() == 0 {
		tape.Running = false
	}
}

func (tape *Tape) Update(tstates int) {
	if tape.pulseWidthCounter > tstates {
		tape.pulseWidthCounter -= tstates
		return
	}

	rest := tstates - tape.pulseWidthCounter

	switch tape.state {
	case TAPE_STATE_BLOCK_START:
		rest = 0
		tape.earBit = false
		tape.blockLength = tape.data.BlockLength(tape.block)

		tape.pulseWidthCounter = LEADER_PULSE_WIDTH
		tape.state = TYPE_STATE_BLOCK_LEADER
		if tape.data.IsHeader(tape.block) {
			tape.pulseCounter = TAPE_HEADER_LEADER_PULSES
		} else {
			tape.pulseCounter = TAPE_DATA_LEADER_PULSES
		}

	case TYPE_STATE_BLOCK_LEADER:
		tape.earBit = !tape.earBit

		tape.pulseCounter--
		if tape.pulseCounter == 0 {
			tape.state = TAPE_STATE_SYNC
			tape.pulseWidthCounter = TAPE_SYNC_HIGH_PULSE_WIDTH
		} else {
			tape.pulseWidthCounter = LEADER_PULSE_WIDTH
		}

	case TAPE_STATE_SYNC:
		tape.earBit = !tape.earBit
		tape.pulseWidthCounter = TAPE_SYNC_LOW_PULSE_WIDTH
		tape.state = TAPE_STATE_BIT_HIGH
		tape.bitMask = 0x80

	case TAPE_STATE_BIT_HIGH:
		tape.earBit = !tape.earBit

		if (tape.data.Byte(tape.position) & tape.bitMask) != 0 {
			tape.pulseWidthCounter = TAPE_BIT_HIGH_PULSE_WIDTH
		} else {
			tape.pulseWidthCounter = TAPE_BIT_LOW_PULSE_WIDTH
		}

		tape.state = TAPE_STATE_BIT_LOW

	case TAPE_STATE_BIT_LOW:
		tape.earBit = !tape.earBit

		if (tape.data.Byte(tape.position) & tape.bitMask) != 0 {
			tape.pulseWidthCounter = TAPE_BIT_HIGH_PULSE_WIDTH
		} else {
			tape.pulseWidthCounter = TAPE_BIT_LOW_PULSE_WIDTH
		}

		tape.bitMask >>= 1

		if tape.bitMask == 0 {
			tape.position++
			tape.blockLength--

			if tape.blockLength == 0 {
				tape.block++
				tape.state = TAPE_STATE_PAUSE

				if tape.block < tape.data.BlockCount() {
					log.Debug("... block %d/%d [%.02f kB]", tape.block, tape.data.BlockCount(), float64(tape.data.BlockLength(tape.block))/1024.0)
				}
			} else {
				tape.bitMask = 0x80
				tape.state = TAPE_STATE_BIT_HIGH
			}
		} else {
			tape.state = TAPE_STATE_BIT_HIGH
		}

	case TAPE_STATE_PAUSE:
		tape.earBit = !tape.earBit

		if tape.block < tape.data.BlockCount() {
			tape.pulseWidthCounter = TAPE_PAUSE_PULSE_WIDTH
			tape.state = TAPE_STATE_PAUSE_END
		} else {
			tape.state = TAPE_STATE_FINISHED
		}

	case TAPE_STATE_PAUSE_END:
		tape.state = TAPE_STATE_BLOCK_START

	case TAPE_STATE_FINISHED:
		elapsedMS := float32(time.Since(tape.start).Milliseconds())
		log.Info("Load finish: %fs", elapsedMS/1000.0)
		tape.Running = false
		//tape.earBit = false
	}

	if tape.pulseWidthCounter >= rest {
		tape.pulseWidthCounter += rest
	}
}

func (tape *Tape) Progress() float32 {
	return float32(tape.position) / float32(tape.data.Length())
}

func (tape *Tape) EarBit() bool {
	return tape.earBit
}
