package synthesizer

import (
	"math"
)

// Lowpass applies a low-pass filter to the frames
// Does not modify the input signal
func Lowpass(fs []float64, freq, delay, sr float64) []float64 {
	output := make([]float64, len(fs))
	copy(output, fs)

	costh := 2. - math.Cos((tau*freq)/sr)
	coef := math.Sqrt(costh*costh-1.) - costh

	for i, a := range output {
		output[i] = a*(1+coef) - delay*coef
		delay = output[i]
	}

	return output
}

// Highpass applies a high-pass filter to the frames.
// Does not modify the input signal
func Highpass(fs []float64, freq, delay, sr float64) []float64 {
	output := make([]float64, len(fs))
	copy(output, fs)

	b := 2. - math.Cos(tau*freq/sr)
	coef := b - math.Sqrt(b*b-1.)

	for i, a := range output {
		output[i] = a*(1.-coef) - delay*coef
		delay = output[i]
	}

	return output
}

// Balance a signal (rescale output signal)
func Balance(signal, comparator, delay []float64, frequency, samplerate float64) []float64 {
	c := make([]float64, len(signal))
	copy(signal, c)

	costh := 2. - math.Cos(tau*frequency/samplerate)
	coef := math.Sqrt(costh*costh-1.) - costh

	for i, s := range signal {
		ss := signal[i]
		if signal[i] < 0 {
			ss = -s
		}
		delay[0] = ss*(1+coef) - (delay[0] * coef)

		if comparator[i] < 0 {
			comparator[i] = -comparator[i]
		}
		delay[1] = comparator[i]*(1+coef) - (delay[1] * coef)
		if delay[0] != 0 {
			c[i] = s * (delay[0] / delay[1])
		} else {
			c[i] = s * delay[1]
		}
	}
	return c
}
