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
	amp      = flag.Float64("a", 1, "Amplitude of signal")
	harms    = flag.Int("h", 1, "Number of harmonics of signal")
	freq     = flag.Float64("f", 440, "Frequency of signal")
	output   = flag.String("o", "", "output file")
)

func main() {
	fmt.Println("Generating triangle wave")
	flag.Parse()

	wfmt := wave.NewWaveFmt(1, 1, 44100, 16, nil)
	triangles := synth.TriangleTable(*harms, *duration*wfmt.SampleRate)
	table := synth.NewGtable(triangles)

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

}
