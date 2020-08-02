package synthesizer

// package to generate oscillators of various shapes

import "math"

const tau = (2 * math.Pi)

type Oscillator struct {
	curfreq  float64
	curphase float64
	incr     float64
	twopiosr float64 // (2*PI) / samplerate
}

// NewOscillator set to a given sample rate
func NewOscillator(sr int) *Oscillator {
	return &Oscillator{
		twopiosr: tau / float64(sr),
	}
}

func triangleCalc(phase float64) float64 {
	val := 2.0*(phase*(1.0/tau)) - 1.0
	if val < 0.0 {
		val = -val
	}
	val = 2.0 * (val - 0.5)
	return val
}

func upwSawtoothCalc(phase float64) float64 {
	val := 2.0*(phase*(1.0/tau)) - 1.0
	return val
}

func downSawtoothCalc(phase float64) float64 {
	val := 1.0 - 2.0*(phase*(1.0/tau))
	return val
}

func squareCalc(phase float64) float64 {
	val := -1.0
	if phase <= math.Pi {
		val = 1.0
	}
	return val
}

func sineCalc(phase float64) float64 {
	return math.Sin(phase)
}

func triangleTick(o *Oscillator, freq float64) float64 {
	if o.curfreq != freq {
		o.curfreq = freq
		o.incr = o.twopiosr * freq
	}
	val := triangleCalc(o.curphase)
	o.curphase += o.incr
	if o.curphase >= tau {
		o.curphase -= tau
	}
	if o.curphase < 0 {
		o.curphase = tau
	}
	return val
}

// A signal that is either 1 or -1 (cycle per math.Pi dist)
func upwardSawtooth(o *Oscillator, freq float64) float64 {
	if o.curfreq != freq {
		o.curfreq = freq
		o.incr = o.twopiosr * freq
	}
	val := upwSawtoothCalc(o.curphase)
	o.curphase += o.incr
	if o.curphase >= tau {
		o.curphase -= tau
	}
	if o.curphase < 0 {
		o.curphase = tau
	}
	return val
}

// A signal that is either 1 or -1 (cycle per math.Pi dist)
func downwardSawtooth(o *Oscillator, freq float64) float64 {
	if o.curfreq != freq {
		o.curfreq = freq
		o.incr = o.twopiosr * freq
	}

	val := downSawtoothCalc(o.curphase)

	o.curphase += o.incr
	if o.curphase >= tau {
		o.curphase -= tau
	}
	if o.curphase < 0 {
		o.curphase = tau
	}
	return val
}

// A signal that is either 1 or -1 (cycle per math.Pi dist)
func sqTick(o *Oscillator, freq float64) float64 {
	if o.curfreq != freq {
		o.curfreq = freq
		o.incr = o.twopiosr * freq
	}

	val := squareCalc(o.curphase)
	o.curphase += o.incr
	if o.curphase >= tau {
		o.curphase -= tau
	}
	if o.curphase < 0 {
		o.curphase = tau
	}
	return val
}

func sineTick(o *Oscillator, freq float64) float64 {
	val := sineCalc(o.curphase)
	if o.curfreq != freq {
		o.curfreq = freq
		o.incr = o.twopiosr * freq
	}

	o.curphase += o.incr
	if o.curphase >= tau {
		o.curphase -= tau
	}
	if o.curphase < 0 {
		o.curphase = tau
	}
	return val

}
