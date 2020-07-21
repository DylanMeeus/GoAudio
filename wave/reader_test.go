package wave

import (
	"testing"
)

var (
	scaleFrameTests = []struct {
		input int
		bits  int
		out   Sample
	}{
		{
			0,
			16,
			0,
		},
		{
			127,
			8,
			1,
		},
		{
			-127,
			8,
			-1,
		},
		{
			32_767,
			16,
			1,
		},
		{
			32_767,
			16,
			1,
		},
		{
			0,
			16,
			0,
		},
		{
			-32_767,
			16,
			-1,
		},
	}
)

func TestScaleFrames(t *testing.T) {
	for _, test := range scaleFrameTests {
		t.Run("", func(t *testing.T) {
			res := scaleFrame(test.input, test.bits)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}
}
