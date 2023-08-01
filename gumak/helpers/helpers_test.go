package helpers

import "testing"

func TestBitCount(t *testing.T) {
	if CountBits(0b0000) != 0 {
		t.Fail()
	}

	if CountBits(0b1) != 1 {
		t.Fail()
	}

	if CountBits(0b1001010) != 3 {
		t.Fail()
	}
}
