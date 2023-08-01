package host

/*
#cgo LDFLAGS: -L../hqx -lhqXscale
#include "../hqx/hqXscale.h"
*/
import "C"

import (
	"fmt"
	"mutex/gumak"
	"mutex/gumak/helpers"
	"path"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// 6144 bytes pixels  (1bit  = 1pixel)
// 768  bytes attribs (1byte = 8x8pixels)
//      BCXXXYYY: B: blink bit, C: contrast bit, XXX: bg color, YYY: fg color

func CallHq3x(input, output []byte, width, height int) {
	C.HQ3X(unsafe.Pointer(&input[0]), C.uint(width), C.uint(height), unsafe.Pointer(&output[0]))
}

func CallHq2x(input, output []byte, width, height int) {
	C.HQ2X(unsafe.Pointer(&input[0]), C.uint(width), C.uint(height), unsafe.Pointer(&output[0]))
}

type helpWin struct {
	show bool

	texture *sdl.Texture
	srcRect sdl.Rect
	dstRect sdl.Rect
}

type Gfx struct {
	innerWidth  int
	innerHeight int

	textureRect       sdl.Rect
	displayTargetRect sdl.Rect

	win      *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	font     *ttf.Font
	fontSize int32

	vramCopy     []byte
	unscaledData []byte

	upscale float64

	fps helpers.Fps

	uTime float32
	start time.Time

	help helpWin

	hq int

	soundMutex *sync.Mutex
}

type Filtering int

const (
	FilteringNearest = Filtering(0)
	FilteringLinear  = Filtering(1)
	FilteringBest    = Filtering(2)
	FilteringHq2x    = Filtering(3)
	FilteringHq3x    = Filtering(4)
)

// Graphics
func (g *Gfx) updateTextureHq(gumak *gumak.Gumak, scale int) {
	data, _, err := g.texture.Lock(&g.textureRect)
	if err != nil {
		panic("Failed to create surface")
	}

	if scale == 3 {
		CallHq3x(g.unscaledData, data, g.innerWidth, g.innerHeight)
	} else {
		CallHq2x(g.unscaledData, data, g.innerWidth, g.innerHeight)
	}

	g.texture.Unlock()
}

func (g *Gfx) updateTexture(gumak *gumak.Gumak) {
	data, pitch, err := g.texture.Lock(&g.textureRect)
	if err != nil {
		panic("Failed to create surface")
	}
	gumak.GetDisplayDataBGR(data, g.vramCopy, pitch)
	g.texture.Unlock()
}

func (g *Gfx) Init(filtering Filtering, innerWidth, innerHeight, outerWidth, outerHeight int, mutex *sync.Mutex) {
	ttf.Init()

	g.uTime = 0
	g.start = time.Now()
	g.soundMutex = mutex
	g.vramCopy = make([]byte, 6912) // TODO: Cleanup.

	g.innerWidth, g.innerHeight = innerWidth, innerHeight

	g.hq = 0
	sdlFiltering := 0
	switch filtering {
	case FilteringNearest:
		sdlFiltering = 0
	case FilteringLinear:
		sdlFiltering = 1
	case FilteringBest:
		sdlFiltering = 2
	case FilteringHq2x:
		g.hq = 1
	case FilteringHq3x:
		g.hq = 2
	}

	textureWidth, textureHeight := int32(innerWidth), int32(innerHeight)
	switch g.hq {
	case 1:
		textureWidth *= 2
		textureHeight *= 2
	case 2:
		textureWidth *= 3
		textureHeight *= 3
	}

	winWidth := g.upscale * float64(outerWidth)
	winHeight := g.upscale * float64(outerHeight)

	w := int32(g.upscale * float64(innerWidth))
	h := int32(g.upscale * float64(innerHeight))
	x := int32((winWidth - float64(w)) / 2.0)
	y := int32((winHeight - float64(h)) / 2.0)
	g.displayTargetRect = sdl.Rect{x, y, w, h}

	g.textureRect = sdl.Rect{0, 0, textureWidth, textureHeight}

	renderer, err := sdl.CreateRenderer(g.win, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		panic(fmt.Sprintf("Failed to create renderer: %s\n", err))
	}

	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, strconv.Itoa(sdlFiltering)) {
		panic(fmt.Sprintf("Failed to set filtering: %s\n", err))
	}

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, textureWidth, textureHeight)
	if err != nil {
		panic(fmt.Sprintf("Failed to create texture: %s\n", err))
	}

	g.renderer = renderer
	g.texture = texture

	if g.hq > 0 {
		g.unscaledData = make([]byte, innerWidth*innerHeight*4)
	}

	g.fontSize = int32(g.upscale * 11.)
	if g.fontSize > 24 {
		g.fontSize = int32(24)
	}

	g.font, err = ttf.OpenFont(path.Join("assets", "Roboto-Black.ttf"), int(g.fontSize))
	if err != nil {
		panic(fmt.Sprintf("Failed to open font: %s", err))
	}

	g.createHelpWin(int32(winWidth), int32(winHeight))
}

func (g *Gfx) Destroy() {
	ttf.Quit()
	g.texture.Destroy()
	g.renderer.Destroy()
}

func (g *Gfx) renderText(text string, color sdl.Color, x, y int32) {
	surface, err := g.font.RenderUTF8Solid(text, color)
	if err != nil {
		panic(fmt.Sprintf("Failed to render font surface: %s", err))
	}

	tex, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(fmt.Sprintf("Failed to create texture from font surface: %s", err))
	}

	rect := sdl.Rect{0, 0, surface.W, surface.H}
	targetRect := sdl.Rect{x, y, surface.W, surface.H}

	g.renderer.Copy(tex, &rect, &targetRect)
}

func (g *Gfx) createHelpWin(w int32, h int32) {
	var err error

	g.help.texture, err = g.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STATIC|sdl.TEXTUREACCESS_TARGET, w, h)
	if err != nil {
		panic(fmt.Sprintf("Failed to create texture: %s\n", err))
	}

	g.help.srcRect = sdl.Rect{0, 0, w, h}
	g.help.dstRect = sdl.Rect{0, 0, w, h}

	// Content.
	g.help.texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	g.renderer.SetRenderTarget(g.help.texture)
	g.renderer.Clear()
	g.renderer.SetDrawColor(0, 0, 0, 128)
	g.renderer.FillRect(&g.displayTargetRect)

	leftCol := int32(g.displayTargetRect.X + g.displayTargetRect.W>>4)
	rightCol := leftCol + int32(g.displayTargetRect.W/2)
	top := int32(g.displayTargetRect.Y + g.displayTargetRect.H>>4)

	g.renderText("F1  - Help", sdl.Color{255, 255, 255, 255}, leftCol, top)
	g.renderText("F2  - Insert & play tape", sdl.Color{255, 0, 0, 255}, leftCol, top+g.fontSize)
	g.renderText("F4  - Reset", sdl.Color{255, 0, 0, 255}, leftCol, top+2*g.fontSize)
	g.renderText("F7  - Toggle audio", sdl.Color{255, 0, 0, 255}, leftCol, top+3*g.fontSize)

	g.renderText("F5  - Quicksave", sdl.Color{255, 255, 0, 255}, rightCol, top)
	g.renderText("F9  - Quickload", sdl.Color{255, 255, 0, 255}, rightCol, top+g.fontSize)
	g.renderText("F11 - Load snapshot", sdl.Color{255, 0, 255, 255}, rightCol, top+2*g.fontSize)
	g.renderText("F12 - Save snapshot", sdl.Color{255, 0, 255, 255}, rightCol, top+3*g.fontSize)

	g.renderer.SetRenderTarget(nil)
}

func (g *Gfx) updateFrame(gumak *gumak.Gumak) {
	g.soundMutex.Lock()
	defer g.soundMutex.Unlock()

	switch g.hq {
	case 0:
		g.updateTexture(gumak)
	case 1:
		gumak.GetDisplayDataBGR(g.unscaledData, g.vramCopy, int(4*g.innerWidth))
	case 2:
		gumak.GetDisplayDataBGR(g.unscaledData, g.vramCopy, int(4*g.innerWidth))
	}
}

func (g *Gfx) Draw(gumak *gumak.Gumak) {
	g.updateFrame(gumak)
	switch g.hq {
	case 1:
		g.updateTextureHq(gumak, 2)
	case 2:
		g.updateTextureHq(gumak, 3)
	}

	g.fps.Update(func(fps float32) {
		g.win.SetTitle(fmt.Sprintf("Gumak Emulator @%.02f FPS [Press F1 for help]", fps))
	})

	g.uTime = float32(time.Since(g.start).Seconds())

	// Clear.
	g.renderer.Clear()

	// Border.
	bgR, bgG, bgB := gumak.GetBackgroundColor()
	g.renderer.SetDrawColor(bgR, bgG, bgB, 255)
	g.renderer.FillRect(&g.displayTargetRect)

	// Screen.
	g.renderer.Copy(g.texture, &g.textureRect, &g.displayTargetRect)

	// Help window.
	if g.help.show {
		g.renderer.Copy(g.help.texture, &g.help.srcRect, &g.help.dstRect)
	}

	// Present.
	g.renderer.Present()
}
