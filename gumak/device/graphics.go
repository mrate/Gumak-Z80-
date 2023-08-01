package device

import "sync"

type Resolution struct {
	W, H int
}

type RGB struct {
	R, G, B uint8
}

var (
	DisplayRes = Resolution{256, 192}
	ScreenRes  = Resolution{320, 240}
	AttribRes  = Resolution{32, 24}
	Colors     = [16]RGB{
		{0, 0, 0},
		{0, 0, 224},
		{224, 0, 0},
		{224, 0, 224},
		{0, 224, 0},
		{0, 224, 224},
		{224, 224, 0},
		{224, 224, 224},

		{0, 0, 0},
		{0, 0, 255},
		{255, 0, 0},
		{255, 0, 255},
		{0, 255, 0},
		{0, 255, 255},
		{255, 255, 0},
		{255, 255, 255},
	}
)

type position struct {
	x, y int
}

func renderPixel(x, y, r, g, b, a int, data []byte, pw int, evenFlash bool, vram []byte) {
	pixel, attr := pixelAttr(vram, x, y)

	blink := (attr & 0b10000000) != 0
	if evenFlash && blink {
		pixel = !pixel
	}

	color := (attr & 0b00111000) >> 3
	if pixel {
		color = attr & 0b00000111
	}

	contr := (attr & 0b01000000) != 0
	if contr {
		color += 8
	}

	pos := (y*int(DisplayRes.W) + x) * pw
	rgb := Colors[color]
	data[pos+r] = rgb.R
	data[pos+g] = rgb.G
	data[pos+b] = rgb.B
	data[pos+a] = 255
}

func pixelAttr(vram []byte, x, y int) (pixel bool, attr uint8) {
	// TODO:

	bit := 7 - (x & 0x7)
	xx := x >> 3
	yy := y >> 3
	w := (DisplayRes.W >> 3)

	// BITMAP:
	// 010 | Y7 Y6 Y2 Y1 Y0 Y5 Y4 Y3 | X4 X3 X2 X1 X0
	y345 := (y & 0b00111000) << 2
	y012 := (y & 0b00000111) << 8
	y67 := (y & 0b11000000) << 5

	idx := xx | y67 | y012 | y345

	pixel = (vram[idx]>>bit)&0x1 != 0
	attr = vram[6144+w*yy+xx]
	return
}

func FillDisplayParallel(data []byte, pitch int, evenFlash bool, vram []byte) {
	// TODO:
	pw := pitch / int(DisplayRes.W)

	var renderQueue = make(chan position)
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range renderQueue {
				// TODO: RGB
				renderPixel(p.x, p.y, 2, 1, 0, 3, data, pw, evenFlash, vram)
			}
		}()
	}

	for y := 0; y < DisplayRes.H; y++ {
		for x := 0; x < DisplayRes.W; x++ {
			renderQueue <- position{x, y}
		}
	}

	close(renderQueue)
	wg.Wait()
}

func FillDisplay(data []byte, pitch int, evenFlash bool, vram []byte, r, g, b, a int) {
	// TODO:
	pw := pitch / int(DisplayRes.W)

	for y := 0; y < DisplayRes.H; y++ {
		for x := 0; x < DisplayRes.W; x++ {
			renderPixel(x, y, r, g, b, a, data, pw, evenFlash, vram)
		}
	}
}
