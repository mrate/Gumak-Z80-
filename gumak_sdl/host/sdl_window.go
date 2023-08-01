package host

import (
	"mutex/gumak"
	"mutex/gumak/log"
	"sync"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
)

type Platform struct {
	upscale            float64
	win                *sdl.Window
	closed             bool
	frameReady         chan bool
	soundGraphicsMutex sync.Mutex

	gfx *Gfx
	snd *Sound
}

// Pixel platform
func createWin(width, height int) *sdl.Window {
	w, h := int32(width), int32(height)

	window, err := sdl.CreateWindow("Gumak Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	return window
}

func (h *Platform) Init(width, height int, upscale float64) {
	h.upscale = upscale
	h.win = createWin(int(float64(width)*upscale), int(float64(height)*upscale))
	h.frameReady = make(chan bool)
}

func (h *Platform) Destroy() {
	h.gfx.Destroy()
	h.snd.Destroy()

	h.win.Destroy()
}

func (h *Platform) WaitFrame() {
	<-h.frameReady
}

func (h *Platform) CreateGfx(filtering Filtering, iw, ih, ow, oh int) *Gfx {
	gfx := new(Gfx)
	gfx.win = h.win
	gfx.upscale = h.upscale
	gfx.Init(filtering, iw, ih, ow, oh, &h.soundGraphicsMutex)

	h.gfx = gfx

	return gfx
}

func (h *Platform) CreateSound(gumak *gumak.Gumak, freq int, samples int) *Sound {
	sound := new(Sound)
	sound.Init(gumak, freq, samples, &h.soundGraphicsMutex, &h.frameReady, h.gfx.vramCopy)
	h.snd = sound

	return sound
}

func (h *Platform) Closed() bool {
	return h.closed
}

func (h *Platform) Close() {
	h.closed = true
}

func (h *Platform) handleEvents(key sdl.Keycode, gumak *gumak.Gumak) {
	switch key {
	case sdl.K_F1:
		h.gfx.help.show = !h.gfx.help.show

	case sdl.K_F2:
		log.Info("Start tape")
		file, err := dialog.File().Filter("TAP tape", "tap").Load()
		if err == nil {
			err := gumak.PlayTape(file)
			if err != nil {
				dialog.Message("Error loading tape: %s", err)
			}
		}

	case sdl.K_F4:
		gumak.Reset()

	case sdl.K_F5:
		gumak.SaveSnapshot("quicksave"+gumak.Model+".z80", nil)

	case sdl.K_F7:
		h.snd.TurnOnOff(!h.snd.IsOn())

	case sdl.K_F9:
		gumak.LoadSnapshot("quicksave"+gumak.Model+".z80", nil)

	case sdl.K_F11:
		file, err := dialog.File().Filter("Snapshot", "sna", "z80").Load()
		if err == nil {
			err = gumak.LoadSnapshot(file, nil)
			if err != nil {
				dialog.Message("Error loading snapshot: %s", err)
				gumak.Reset()
			}
		}

	case sdl.K_F12:
		file, err := dialog.File().Filter("Snapshot", "sna", "z80").Save()
		if err == nil {
			err = gumak.SaveSnapshot(file, nil)
			if err != nil {
				dialog.Message("Error saving snapshot: %s", err)
			}
		}
	}
}

func (h *Platform) Update(gumak *gumak.Gumak) {
	h.soundGraphicsMutex.Lock()
	defer h.soundGraphicsMutex.Unlock()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			h.closed = true
		case *sdl.KeyboardEvent:
			keys, ok := KeyMap[e.Keysym.Sym]
			if ok {
				for _, k := range keys {
					gumak.HandleKey(k, e.Type == sdl.KEYDOWN)
				}
			} else if e.Type == sdl.KEYDOWN {
				h.handleEvents(e.Keysym.Sym, gumak)
			}
		}
	}
}
