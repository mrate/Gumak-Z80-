package device

type Beeper struct {
	value  bool
	filter filter
}

type filter struct {
	reqValue bool
	value    float64
	inc      float64
}

func (b *Beeper) Init(filt float64) {
	b.filter.inc = filt
}

func (b *Beeper) Beep(val bool) {
	b.value = val
}

func (b *Beeper) Reset() {
	b.value = false
}

func (b *Beeper) Sample() float64 {
	b.filter.Update(b.value)
	return b.filter.value
}

func (f *filter) Update(value bool) {
	f.reqValue = value

	if f.reqValue {
		if f.value < 0.5-f.inc {
			f.value += f.inc
		} else {
			f.value = 0.5
		}
	} else {
		if f.value > -0.5+f.inc {
			f.value -= f.inc
		} else {
			f.value = -0.5
		}
	}
}
