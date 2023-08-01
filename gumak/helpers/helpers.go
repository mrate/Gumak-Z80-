package helpers

func Exchange8(a *uint8, b *uint8) {
	tmp := *a
	*a = *b
	*b = tmp
}

func Exchange16(a *uint16, b *uint16) {
	tmp := *a
	*a = *b
	*b = tmp
}

func CountBits(value int) int {
	bits := 0
	for value > 0 {
		bits++
		value = value & (value - 1)
	}

	return bits
}

func To16(l uint8, h uint8) uint16 {
	return uint16(h)<<8 | uint16(l)
}

func To8(v uint16) (l, h uint8) {
	l = uint8(v & 0xff)
	h = uint8((v >> 8) & 0xff)
	return
}

func Approxsin(t float64) float64 {
	j := t * 0.15915
	j = j - float64(int(j))
	return 20.785 * j * (j - 0.5) * (j - 1.0)
}
