package wave

import (
	"testing"
)

var (
	rescaleFrameTests = []struct {
		input Frame
		bits  int
		out   int
	}{
		{
			Frame(1),
			8,
			127,
		},
		{
			Frame(-1),
			8,
			-127,
		},
		{
			Frame(1),
			16,
			32_767,
		},
		{
			Frame(-1),
			16,
			-32_767,
		},
	}
)

func TestRescaleFrames(t *testing.T) {
	for _, test := range rescaleFrameTests {
		t.Run("", func(t *testing.T) {
			res := rescaleFrame(test.input, test.bits)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}

}
