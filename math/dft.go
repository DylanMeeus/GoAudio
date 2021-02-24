package math

import (
	"github.com/DylanMeeus/GoAudio/wave"
	"math"
	"math/cmplx"
)

const (
	tau = 2. * math.Pi
)

// DFT is a discrete fourier transformation on the input frames
// DEPRECATED
// Please use FFT unless you are sure  you want this one..
func DFT(input []wave.Frame) []complex128 {
	N := len(input)

	output := make([]complex128, len(input))

	reals := make([]float64, len(input))
	imgs := make([]float64, len(input))
	for i, frame := range input {
		for n := 0; n < N; n++ {
			reals[i] += float64(frame) * math.Cos(float64(i*n)*tau/float64(N))
			imgs[i] += float64(frame) * math.Sin(float64(i*n)*tau/float64(N))
		}

		reals[i] /= float64(N)
		imgs[i] /= float64(N)
	}

	for i := 0; i < len(reals); i++ {
		output[i] = complex(reals[i], imgs[i])
	}

	return output
}

// HFFT mutates freqs!
func HFFT(input []wave.Frame, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(input[0], 0)
		return
	}

	h := n / 2

	HFFT(input, freqs, h, 2*step)
	HFFT(input[step:], freqs[step:], h, 2*step)

	for k := 0; k < h; k++ {
		a := -2 * math.Pi * float64(k) * float64(n)
		e := cmplx.Rect(1, a) * freqs[k+h]
		freqs[k], freqs[k+h] = freqs[k]+e, freqs[k]-e
	}
}

// FFT (Fast Fourier Transform) implementation
func FFT(input []wave.Frame) []complex128 {
	freqs := make([]complex128, len(input))
	HFFT(input, freqs, len(input), 1)
	return freqs
}
