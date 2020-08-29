package main

import (
	"flag"
	"fmt"
	synth "github.com/DylanMeeus/GoAudio/synthesizer"
	"github.com/DylanMeeus/GoAudio/wave"
)

// example use of the oscillator to generate different waveforms
var (
	duration = flag.Int("d", 10, "duration of signal in seconds")
	shape    = flag.String("s", "square", "waveshape to generate")
	amp      = flag.Float64("a", 1, "Amplitude of signal")
	harms    = flag.Int("h", 1, "Number of harmonics of signal")
	freq     = flag.Float64("f", 440, "Frequency of signal")
	output   = flag.String("o", "", "output file")
)

var (
	shapefunc = map[string]func(int, int) []float64{
		"triangle": synth.TriangleTable,
		"square":   synth.SquareTable,
		"saw":      synth.SawTable,
	}
)

func main() {
	flag.Parse()

	wfmt := wave.NewWaveFmt(1, 1, 44100, 16, nil)

	waveform := shapefunc[*shape]
	if waveform == nil {
		waveform = synth.SquareTable
	}
	table := synth.NewGtable(waveform(*harms, *duration*wfmt.SampleRate))

	osc, err := synth.NewLookupOscillator(wfmt.SampleRate, table, 0.0)
	if err != nil {
		panic(err)
	}

	samples := osc.BatchTruncateTick(*freq, *duration*wfmt.SampleRate)
	frames := []wave.Frame{}
	for _, s := range samples {
		frames = append(frames, wave.Frame(s))
	}
	if err := wave.WriteFrames(frames, wfmt, *output); err != nil {
		panic(err)
	}
	fmt.Println("done")
}
