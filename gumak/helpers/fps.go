package helpers

import "time"

type Fps struct {
	frames   uint64
	lastTime int64
}

func (f *Fps) Update(fps func(float32)) {
	f.frames++

	time := time.Now().UnixMilli()
	elapsed := time - f.lastTime
	if elapsed >= 1000 {
		val := float32(f.frames) * 1000.0 / float32(elapsed)
		fps(val)

		f.lastTime = time
		f.frames = 0
	}
}
