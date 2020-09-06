package synthesizer_test

import (
	synth "github.com/DylanMeeus/GoAudio/synthesizer"

	"testing"
)

var (
	lowpassTests = []struct {
		input []float64
	}{
		{
			input: []float64{},
		},
		{
			input: []float64{1, 2, 3},
		},
	}

	highpassTests = []struct {
		input []float64
	}{
		{
			input: []float64{},
		},
		{
			input: []float64{1, 2, 3},
		},
	}
)

func TestLowpassFilter(t *testing.T) {
	// Test to make sure  that the function is non-destructive to the input
	for _, test := range lowpassTests {
		t.Run("", func(t *testing.T) {
			c := make([]float64, len(test.input))
			copy(test.input, c)
			_ = synth.Lowpass(test.input, 440, 1, 44100)
			// make sure that the input is not modified
			if !floatsEqual(test.input, c) {
				t.Fatal("Function modified source!")
			}
		})
	}
}

func TestHighpassFilter(t *testing.T) {
	// Test to make sure  that the function is non-destructive to the input
	for _, test := range highpassTests {
		t.Run("", func(t *testing.T) {
			c := make([]float64, len(test.input))
			copy(test.input, c)
			_ = synth.Highpass(test.input, 440, 1, 44100)
			// make sure that the input is not modified
			if !floatsEqual(test.input, c) {
				t.Fatal("Function modified source!")
			}
		})
	}
}
func floatsEqual(fs1, fs2 []float64) bool {
	if len(fs1) != len(fs2) {
		return false
	}

	for i := range fs1 {
		if fs2[i] != fs2[i] {
			return false
		}
	}
	return true
}
