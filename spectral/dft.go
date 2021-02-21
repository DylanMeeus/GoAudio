package spectral

import (
	"github.com/DylanMeeus/GoAudio/wave"
	"math"
)

const (
	tau = 2. * math.Pi
)

// DFT is a discrete fourier transformation on the input frames
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
