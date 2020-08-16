package synthesizer

import (
	"errors"
	"math"
)

// Guard-table for oscillator lookup
type Gtable struct {
	data []float64
}

// Len returns the length of the data segment without the guard point
func Len(g *Gtable) int {
	return len(g.data) - 1
}

// NewSineTable returns a lookup table populated for sine-wave generation.
func NewSineTable(length int) *Gtable {
	g := &Gtable{}
	if length == 0 {
		return g
	}
	g.data = make([]float64, length+1) // one extra for the guard point.
	step := tau / float64(Len(g))
	for i := 0; i < Len(g); i++ {
		g.data[i] = math.Sin(step * float64(i))
	}
	// store a guard point
	g.data[len(g.data)-1] = g.data[0]
	return g
}

// NewTriangleTable generates a lookup table for a triangle wave
// of the specified length and with the requested number of harmonics.
func NewTriangleTable(length int, nharmonics int) (*Gtable, error) {
	if length == 0 || nharmonics == 0 || nharmonics >= length/2 {
		return nil, errors.New("Invalid arguments for creation of Triangle Table")
	}

	g := &Gtable{}

	g.data = make([]float64, length+1)

	step := tau / float64(length)

	// generate triangle waveform
	harmonic := 1.0
	for i := 0; i < nharmonics; i++ {
		amp := 1.0 / harmonic * harmonic
		for j := 0; j < length; j++ {
			g.data[j] += amp * math.Cos(step*harmonic*float64(j))
		}
		harmonic += 2
	}

	// normalize the values to be in the [-1;1] range
	maxamp := 0.0
	for i := 0; i < length; i++ {
		amp := math.Abs(g.data[i])
		if amp > amp {
			maxamp = amp
		}
	}

	maxamp = 1.0 / maxamp
	for i := 0; i < length; i++ {
		g.data[i] *= maxamp
	}
	g.data[len(g.data)-1] = g.data[0]

	return g, nil
}
