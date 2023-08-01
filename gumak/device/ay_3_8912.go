package device

import (
	"math"
	"mutex/gumak/helpers"
)

type Channel struct {
	freq       float64
	onOff      bool
	volume     uint8 // 0 - 15
	autoVolume bool
	noiseOn    bool
}

type Envelope struct {
	freq  float64
	shape uint8
	time  float64
	cycle float64
}

type AY_3_8912 struct {
	selectedReg uint8

	regs [15]uint8

	noiseFreq float64

	channelA Channel
	channelB Channel
	channelC Channel

	env Envelope

	time float64

	dc      Dc
	dcIndex int
}

/*
// D/C table, taken from AY-3-8912 manual.
var Volumes = []float64{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0.125,
	0.1515,
	0.25,
	0.303,
	0.5,
	0.707,
	1,
}
*/

// D/C table.
// src: https://github.com/true-grue/ayumi/blob/master/ayumi.c
// Manual contains only couple of values...
var Volumes = []float64{
	0.0,
	0.00999465934234,
	0.0144502937362,
	0.0210574502174,
	0.0307011520562,
	0.0455481803616,
	0.0644998855573,
	0.107362478065,
	0.126588845655,
	0.20498970016,
	0.292210269322,
	0.372838941024,
	0.492530708782,
	0.635324635691,
	0.805584802014,
	1.0,
}

// Envelope

func (e *Envelope) Update(timeDelta float64) uint8 {

	e.time += timeDelta
	cycle, fract := math.Modf(e.freq * e.time)

	if cycle > 0 {
		// Cont.
		if (e.shape & 0b0001) == 0 {
			return 0
		}

		// Hold
		if (e.shape & 0b1000) != 0 {
			b1, b2 := (e.shape&0b0010)>>1, (e.shape&0b0100)>>2
			if (b1 ^ b2) == 0 {
				return 0
			} else {
				return 1
			}
		}
	}

	// Alternate
	if (e.shape&0b0100) != 0 && cycle != e.cycle {
		e.shape ^= 0b0010
	}

	e.cycle = cycle

	// Dir (1 up, 0 down)
	if (e.shape & 0b0010) == 0 {
		fract = 0.999999 - fract
	}

	return uint8(16 * fract)
}

func (e *Envelope) SetFreq(freq uint16) {
	e.time = 0
	e.cycle = 0
	e.freq = (1789772.5 / (256 * float64(freq)))
}

func (e *Envelope) SetShape(shape uint8) {
	e.time = 0
	e.cycle = 0
	e.shape = shape
}

// Channel
func (t *Channel) Volume(env uint8) float64 {
	if t.autoVolume {
		return Volumes[env]
	} else {
		return Volumes[t.volume]
	}
}

func (t *Channel) Noise(noise float64) float64 {
	if t.noiseOn {
		return noise
	} else {
		return 0
	}
}

func (t *Channel) Tone(time float64) float64 {
	if !t.onOff {
		return 0
	}

	// Square wave from sine harmonics.
	// https://www.youtube.com/watch?v=72dI7dB3ZvQ&t=1s
	p := math.Pi

	freqTime := t.freq * 2. * math.Pi * time

	a := 0.0
	b := 0.0
	c := 0.0
	for n := 1.; n < 20; n++ {
		c += freqTime
		a += -helpers.Approxsin(c) / n
		b += -helpers.Approxsin(c-p*n) / n
	}

	return (2 / math.Pi) * (a - b)
}

func (t *Channel) SetFreq(freq uint16) {
	if freq == 0 {
		freq = 1
	}

	t.freq = 1789772.5 / (16 * float64(freq))
}

func (t *Channel) SetVolume(volume uint8) {
	t.autoVolume = (volume & 0b10000) != 0
	t.volume = volume & 0b1111
}

// Sound chip
func (a *AY_3_8912) SetNoiseFreq(freq uint8) {
	if freq == 0 {
		freq = 1
	}
	a.noiseFreq = 1789772.5 / (16 * float64(freq))
}

func (a *AY_3_8912) SelectRegister(reg uint8) {
	if reg <= 14 {
		a.selectedReg = reg
	} else {
		a.selectedReg = 14
	}
}

func (a *AY_3_8912) Write(val uint8) {
	a.regs[a.selectedReg] = val

	switch a.selectedReg {
	case 0, 1:
		a.channelA.SetFreq(helpers.To16(a.regs[0], a.regs[1]&0b1111))
	case 2, 3:
		a.channelB.SetFreq(helpers.To16(a.regs[2], a.regs[3]&0b1111))
	case 4, 5:
		a.channelC.SetFreq(helpers.To16(a.regs[4], a.regs[5]&0b1111))
	case 6:
		a.SetNoiseFreq(a.regs[6] & 0b11111)
	case 7:
		a.channelA.onOff = (a.regs[7] & 0b001) == 0
		a.channelB.onOff = (a.regs[7] & 0b010) == 0
		a.channelC.onOff = (a.regs[7] & 0b100) == 0
		a.channelA.noiseOn = (a.regs[7] & 0b001000) == 0
		a.channelB.noiseOn = (a.regs[7] & 0b010000) == 0
		a.channelC.noiseOn = (a.regs[7] & 0b100000) == 0
	case 8:
		a.channelA.SetVolume(a.regs[8] & 0b1111)
	case 9:
		a.channelB.SetVolume(a.regs[9] & 0b1111)
	case 10:
		a.channelC.SetVolume(a.regs[10] & 0b1111)
	case 11, 12:
		a.env.SetFreq(helpers.To16(a.regs[11], a.regs[12]))
	case 13:
		a.env.SetShape(a.regs[13] & 0b1111)
	}
}

func (a *AY_3_8912) Noise(time float64) float64 {
	// TODO: Pseudo-random noise.
	freqTime := a.noiseFreq * 2. * math.Pi * time

	v := 1000.
	for i := 1.; i < 20; i++ {
		v *= helpers.Approxsin(i * freqTime / (i + 1))
	}

	return v
}

func (a *AY_3_8912) Read() uint8 {
	return a.regs[a.selectedReg]
}

func (a *AY_3_8912) Sample(timeDelta float64) float64 {
	a.time += timeDelta

	coef := 0.33
	env := a.env.Update(timeDelta)
	noise := a.Noise(a.time)

	value := coef * a.channelA.Volume(env) * (a.channelA.Tone(a.time) + a.channelA.Noise(noise))
	value += coef * a.channelB.Volume(env) * (a.channelB.Tone(a.time) + a.channelB.Noise(noise))
	value += coef * a.channelC.Volume(env) * (a.channelC.Tone(a.time) + a.channelC.Noise(noise))

	value = a.dc.Filter(a.dcIndex, value)
	a.dcIndex = (a.dcIndex + 1) & (DC_FILTER_SIZE - 1)

	return value
}
