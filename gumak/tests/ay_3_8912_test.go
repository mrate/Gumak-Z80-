package tests

import (
	"mutex/gumak/device"
	"testing"
)

func equals(a, b, eps float64) bool {
	return a >= b-eps && a <= b+eps
}

func testShape(t *testing.T, name string, s int, results []uint8) {
	var e device.Envelope

	e.SetFreq(6991) // ~1

	e.SetShape(uint8(s))

	val := e.Update(0)

	// fmt.Printf("Shape: %d [%s]\n", s, name)
	// fmt.Printf("%d,\n", val)

	if results[0] != val {
		t.Fatalf("Expected %d got %d", results[0], val)
	}

	for i, time := 0, 0.; i < 14; i, time = i+1, time+0.2 {
		val = e.Update(0.2)
		// fmt.Printf("%d,\n", val)

		if results[i+1] != val {
			t.Fatalf("Expected %d got %d", results[i+1], val)
		}
	}

	// fmt.Println("-------------")
}

func TestEnvelope(t *testing.T) {

	{
		results := []uint8{
			15,
			12,
			9,
			6,
			3,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		}

		testShape(t, `\______`, 0b0000, results)
		testShape(t, `\______`, 0b0100, results)
		testShape(t, `\______`, 0b1000, results)
		testShape(t, `\______`, 0b1100, results)
	}

	{
		results := []uint8{
			0,
			3,
			6,
			9,
			12,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		}

		testShape(t, `/|_____`, 0b0010, results)
		testShape(t, `/|_____`, 0b0110, results)
		testShape(t, `/|_____`, 0b1010, results)
		testShape(t, `/|_____`, 0b1110, results)
	}

	{
		results := []uint8{
			15,
			12,
			9,
			6,
			3,
			15,
			12,
			9,
			6,
			3,
			15,
			12,
			9,
			6,
			3,
		}

		testShape(t, `|\|\|\|\`, 0b0001, results)
	}

	{
		results := []uint8{
			15,
			12,
			9,
			6,
			3,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		}

		testShape(t, `\_______`, 0b1001, results)
	}

	{
		results := []uint8{
			15,
			12,
			9,
			6,
			3,
			0,
			3,
			6,
			9,
			12,
			15,
			12,
			9,
			6,
			3,
		}

		testShape(t, `\/\/\/\/\/\`, 0b0101, results)
	}

	{
		results := []uint8{
			15,
			12,
			9,
			6,
			3,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
		}

		testShape(t, `\|^^^^^^^^^`, 0b1101, results)
	}

	{
		results := []uint8{
			0,
			3,
			6,
			9,
			12,
			0,
			3,
			6,
			9,
			12,
			0,
			3,
			6,
			9,
			12,
		}

		testShape(t, `/|/|/|/|/|`, 0b0011, results)
	}

	{
		results := []uint8{
			0,
			3,
			6,
			9,
			12,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
		}

		testShape(t, `/^^^^^^`, 0b1011, results)
	}

	{
		results := []uint8{
			0,
			3,
			6,
			9,
			12,
			15,
			12,
			9,
			6,
			3,
			0,
			3,
			6,
			9,
			12,
		}

		testShape(t, `/\/\/\/\/\/`, 0b0111, results)
	}

	{
		results := []uint8{
			0,
			3,
			6,
			9,
			12,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		}

		testShape(t, `/|_______`, 0b1111, results)
	}
}
