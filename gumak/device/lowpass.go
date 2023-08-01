package device

import "math"

type Lowpass struct {
	cutoff float64
	ePow   float64
	output float64
}

func (l *Lowpass) Init(cutoff, deltaTime float64) {
	l.cutoff = cutoff
	l.ePow = 1 - math.Exp(-deltaTime*2*math.Pi*cutoff)
	l.output = 0
}

func (l *Lowpass) Update(input float64) float64 {
	l.output += (input - l.output) * l.ePow
	return l.output
}

const (
	DC_FILTER_SIZE = 1024
)

type Dc struct {
	sum   float64
	delay [DC_FILTER_SIZE]float64
}

func (dc *Dc) Filter(index int, x float64) float64 {
	dc.sum += -dc.delay[index] + x
	dc.delay[index] = x
	return x - dc.sum/float64(DC_FILTER_SIZE)
}
