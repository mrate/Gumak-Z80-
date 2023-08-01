package gumak

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"os"

	"mutex/gumak/device"
	"mutex/gumak/formats"
	"mutex/gumak/log"
	"mutex/gumak/z80"
)

//go:embed roms/*.rom
var embedContent embed.FS

type Key int

const (
	KeyShift = iota
	KeyZ
	KeyX
	KeyC
	KeyV
	KeyA
	KeyS
	KeyD
	KeyF
	KeyG
	KeyQ
	KeyW
	KeyE
	KeyR
	KeyT
	Key1
	Key2
	Key3
	Key4
	Key5
	Key0
	Key9
	Key8
	Key7
	Key6
	KeyP
	KeyO
	KeyI
	KeyU
	KeyY
	KeyEnter
	KeyL
	KeyK
	KeyJ
	KeyH
	KeySpace
	KeySymbolShift
	KeyM
	KeyN
	KeyB
)

type config struct {
	frequency       int
	tStatesPerFrame int
	roms            []string
	paging          bool
}

var machines = map[string]config{
	"48": {
		frequency:       3500000,
		tStatesPerFrame: 224 * 312,
		roms:            []string{"48.rom"},
		paging:          false,
	},
	"128": {
		frequency:       3546900,
		tStatesPerFrame: 228 * 311,
		roms:            []string{"128-0.rom", "128-1.rom"},
		paging:          true,
	},
}

type Gumak struct {
	Model string

	Cpu       *z80.CPU
	Ram       *device.Ram
	Ula       *device.Ula
	Beeper    *device.Beeper
	Ay_3_8912 *device.AY_3_8912

	roms    []string
	romPath string

	sampleCounter float64 //
	sampleTime    float64 // Time of single audio frame in seconds.

	tStatesFrame   int
	tStatesSeconds float64

	tapeLoading  bool
	tapeFinished chan bool

	lowPass device.Lowpass
}

// Main

func CreateNew(is48K bool, audioFreq int) (*Gumak, error) {
	log.SetLevel(log.LevelDebug)
	log.Info("=== Gumak started ===")

	machineSettings := machines["128"]
	machine := "128"
	if is48K {
		machineSettings = machines["48"]
		machine = "48"
	}

	log.TraceEnable(0)
	log.Trace(2, "=== TRACE ENABLED ===\n")

	// HW
	cpu := new(z80.CPU)
	cpu.Init(machineSettings.frequency, machineSettings.tStatesPerFrame, nil)

	ram := new(device.Ram)
	ram.Init()

	beeper := new(device.Beeper)
	beeper.Init(0.25)

	ay_3_8192 := new(device.AY_3_8912)

	ula := new(device.Ula)
	ula.Init(ram, beeper, ay_3_8192, machineSettings.paging)

	gumak := new(Gumak)
	gumak.Cpu = cpu
	gumak.Ram = ram
	gumak.Ula = ula
	gumak.Beeper = beeper
	gumak.Ay_3_8912 = ay_3_8192
	gumak.Model = machine
	gumak.roms = machineSettings.roms
	gumak.romPath = "roms"
	gumak.sampleTime = 1 / float64(audioFreq)

	gumak.lowPass.Init(100, gumak.sampleTime)

	// Memory + IO bus.
	cpu.Pin.Bus = func() {
		switch {
		case cpu.Pin.MREQ: // Memory request
			if cpu.Pin.RD {
				cpu.Pin.DATA = ram.Read(cpu.Pin.ADDR)
			} else if cpu.Pin.WR {
				ram.Write(cpu.Pin.ADDR, cpu.Pin.DATA)
			}
		case cpu.Pin.IOREQ: // I/O request
			if cpu.Pin.RD {
				cpu.Pin.DATA = ula.Read(cpu.Pin.ADDR)
			} else {
				ula.Write(cpu.Pin.ADDR, cpu.Pin.DATA)
			}
		}
	}

	gumak.tStatesSeconds = cpu.TStateUs / 1e6
	gumak.tapeFinished = make(chan bool, 16)

	// Run
	err := gumak.Reset()
	if err != nil {
		return nil, err
	}

	return gumak, nil
}

func (g *Gumak) Reset() error {
	g.Cpu.Restart()
	g.Ram.Init()

	g.tStatesFrame = 0
	g.sampleCounter = 0

	g.Beeper.Reset()

	for i, rom := range g.roms {
		_, err := g.Ram.LoadRom(embedContent, fmt.Sprintf("%s/%s", g.romPath, rom), i, 0x4000)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Gumak) Tick() bool {
	if g.Ula.Tape.Running && !g.tapeLoading {
		g.tapeLoading = true

		go func() {
			for g.Ula.Tape.Running {
				t := g.Cpu.Tick()
				g.Ula.Tape.Update(t)
			}

			g.tapeFinished <- true
		}()
		return true
	}

	if g.tapeLoading {
		select {
		case <-g.tapeFinished:
			g.tapeLoading = false
		default:
			return true
		}
	}

	if g.tStatesFrame < g.Cpu.TStatesPerFrame {
		t := g.Cpu.Tick()

		g.tStatesFrame += t
		g.sampleCounter += float64(t) * g.tStatesSeconds

		g.Cpu.Pin.INT = false
	}

	if g.tStatesFrame >= g.Cpu.TStatesPerFrame {
		g.Ula.UpdateEndFrame()
		g.tStatesFrame -= g.Cpu.TStatesPerFrame
		g.Cpu.Pin.INT = true
		return true
	}

	return false
}

// Graphics

func InnerResolution() (width int, height int) {
	return device.DisplayRes.W, device.DisplayRes.H
}

func OuterResolution() (width int, height int) {
	return device.ScreenRes.W, device.ScreenRes.H
}

func (g *Gumak) CopyVRam(output []byte) {
	// TODO: Improve.
	// Since we are syncing by audio we need to make copy of a VRAM
	// exactly when frame should be rendered (not later).
	v := g.Ram.Bank(g.Ula.VRamBank)

	for i := range output {
		output[i] = v[i]
	}
}

func (g *Gumak) GetDisplayDataRGB(data, vram []byte, pitch int) {
	device.FillDisplay(data, pitch, g.Ula.Frames >= 16, vram, 0, 1, 2, 3)
}

func (g *Gumak) GetDisplayDataBGR(data, vram []byte, pitch int) {
	device.FillDisplay(data, pitch, g.Ula.Frames >= 16, vram, 2, 1, 0, 3)
}

func (g *Gumak) GetDisplayDataRGBVram(data []byte, pitch int) {
	device.FillDisplay(data, pitch, g.Ula.Frames >= 16, g.Ram.Bank(g.Ula.VRamBank), 0, 1, 2, 3)
}

func (g *Gumak) GetDisplayDataBGRVram(data []byte, pitch int) {
	device.FillDisplay(data, pitch, g.Ula.Frames >= 16, g.Ram.Bank(g.Ula.VRamBank), 2, 1, 0, 3)
}

func (gum *Gumak) GetBackgroundColor() (r uint8, g uint8, b uint8) {
	color := device.Colors[gum.Ula.BorderColor]
	return color.R, color.G, color.B
}

// Input

func (g *Gumak) HandleKey(key Key, down bool) {
	rowIndex := int(key) / 5
	bit := uint8(1 << (int(key) % 5))

	if down {
		g.Ula.Keyboard[rowIndex] &= ^bit
	} else {
		g.Ula.Keyboard[rowIndex] |= bit
	}
}

// I/O

func (g *Gumak) PlayTape(file string) error {
	source, err := formats.NewTape(file)
	if err != nil {
		return err
	}

	err = g.Ula.Tape.Load(file, source)
	if err != nil {
		return err
	}

	g.Ula.Tape.Start()
	return nil
}

func (g *Gumak) LoadSnapshot(filename string, reader io.Reader) error {
	snapshot, err := formats.NewSnapshot(filename)
	if err != nil {
		return err
	}

	if reader == nil {
		reader, err = os.Open(filename)
		if err != nil {
			return err
		}
	}

	return snapshot.Load(reader, g.Cpu, g.Ula, g.Ram)
}

func (g *Gumak) SaveSnapshot(filename string, writer io.Writer) error {
	snapshot, err := formats.NewSnapshot(filename)
	if err != nil {
		return err
	}

	if writer == nil {
		writer, err = os.Create(filename)
		if err != nil {
			return err
		}
	}

	return snapshot.Save(writer, g.Cpu, g.Ula, g.Ram)
}

// Audio
func (g *Gumak) AudioSampleReady() bool {
	return g.sampleCounter >= g.sampleTime
}

func (g *Gumak) PopAudioSample() uint8 {
	g.sampleCounter -= g.sampleTime

	beep := g.Beeper.Sample()
	sound := g.Ay_3_8912.Sample(g.sampleTime)

	mix := 0.5*beep + 0.5*sound
	return uint8(128 + 128*mix)
}
