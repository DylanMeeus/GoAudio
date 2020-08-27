package synthesizer

import (
	"errors"
)

// LookupOscillator is an oscillator that's more gentle on your CPU
// By performing a table lookup to generate the required waveform..
type LookupOscillator struct {
	Oscillator
	Table      *Gtable
	SizeOverSr float64 // convenience variable for calculations
}

// NewLookupOscillator creates a new oscillator which
// performs a table-lookup to generate the required waveform
func NewLookupOscillator(sr int, t *Gtable, phase float64) (*LookupOscillator, error) {
	if t == nil || len(t.data) == 0 {
		return nil, errors.New("Invalid table provided for lookup oscillator")
	}

	return &LookupOscillator{
		Oscillator: Oscillator{
			curfreq:  0.0,
			curphase: float64(Len(t)) * phase,
			incr:     0.0,
		},
		Table:      t,
		SizeOverSr: float64(Len(t)) / float64(sr),
	}, nil

}

// TruncateTick performs a lookup and truncates the value
// index down (if the index for lookup = 10.5, return index 10)
func (l *LookupOscillator) TruncateTick(freq float64) float64 {
	return l.BatchTruncateTick(freq, 1)[0]
}

// BatchTruncateTick returns a slice of samples from the oscillator of the requested length
func (l *LookupOscillator) BatchTruncateTick(freq float64, nframes int) []float64 {
	out := make([]float64, nframes)
	for i := 0; i < nframes; i++ {
		index := l.curphase
		if l.curfreq != freq {
			l.curfreq = freq
			l.incr = l.SizeOverSr * l.curfreq
		}
		curphase := l.curphase
		curphase += l.incr
		for curphase > float64(Len(l.Table)) {
			curphase -= float64(Len(l.Table))
		}
		for curphase < 0.0 {
			curphase += float64(Len(l.Table))
		}
		l.curphase = curphase
		out[i] = l.Table.data[int(index)]
	}
	return out
}

// InterpolateTick performs a lookup but interpolates the value if the
// requested index does not appear in the table.
func (l *LookupOscillator) InterpolateTick(freq float64) float64 {
	return l.BatchInterpolateTick(freq, 1)[0]
}

// BatchInterpolateTick performs a lookup for N frames, and interpolates the value if the
// requested index does not appear in the table.
func (l *LookupOscillator) BatchInterpolateTick(freq float64, nframes int) []float64 {
	out := make([]float64, nframes)
	for i := 0; i < nframes; i++ {
		baseIndex := int(l.curphase)
		nextIndex := baseIndex + 1
		if l.curfreq != freq {
			l.curfreq = freq
			l.incr = l.SizeOverSr * l.curfreq
		}
		curphase := l.curphase
		frac := curphase - float64(baseIndex)
		val := l.Table.data[baseIndex]
		slope := l.Table.data[nextIndex] - val
		val += frac * slope
		curphase += l.incr

		for curphase > float64(Len(l.Table)) {
			curphase -= float64(Len(l.Table))
		}
		for curphase < 0.0 {
			curphase += float64(Len(l.Table))
		}

		l.curphase = curphase
		out[i] = val
	}
	return out
}
