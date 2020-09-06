package synthesizer

import (
	"math"
)

// Lowpass applies a low-pass filter to the frames
// Does not modify the input signal
func Lowpass(fs []float64, freq, delay, sr float64) []float64 {
	output := make([]float64, len(fs))
	copy(fs, output)

	b := 2. - math.Cos(tau*freq/sr)
	coef := math.Sqrt(b*b-1.) - b

	for i, a := range output {
		output[i] = a*(1+coef) - delay*coef
		delay = a
	}

	return output
}

// Highpass applies a high-pass filter to the frames.
// Does not modify the input signal
func Highpass(fs []float64, freq, delay, sr float64) []float64 {
	output := make([]float64, len(fs))
	copy(fs, output)

	b := 2. - math.Cos(tau*freq/sr)
	coef := b - math.Sqrt(b*b-1.)

	for i, a := range output {
		output[i] = a*(1.-coef) - delay*coef
		delay = a
	}

	return output
}
