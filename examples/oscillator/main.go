package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DylanMeeus/GoAudio/breakpoint"
	synth "github.com/DylanMeeus/GoAudio/synthesizer"
	"github.com/DylanMeeus/GoAudio/wave"
	"io/ioutil"
)

var stringToShape = map[string]synth.Shape{
	"sine":     0,
	"square":   1,
	"downsaw":  2,
	"upsaw":    3,
	"triangle": 4,
}

// example use of the oscillator to generate different waveforms
var (
	duration   = flag.Int("d", 10, "duration of signal")
	shape      = flag.String("s", "sine", "One of: sine, square, triangle, downsaw, upsaw")
	amppoints  = flag.String("a", "", "amplitude breakpoints file")
	freqpoints = flag.String("f", "", "frequency breakpoints file")
	output     = flag.String("o", "", "output file")
)

// Generate an oscillator of a given shape and duration
// Modified by amplitude / frequency breakpoints
func main() {
	flag.Parse()
	fmt.Println("usage: go run main -d {dur} -s {shape} -a {amps} -f {freqs} -o {output}")
	if output == nil {
		panic("please provide an output file")
	}

	wfmt := wave.NewWaveFmt(1, 1, 44100, 16, nil)
	amps, err := ioutil.ReadFile(*amppoints)
	if err != nil {
		panic(err)
	}
	ampPoints, err := breakpoint.ParseBreakpoints(bytes.NewReader(amps))
	if err != nil {
		panic(err)
	}
	ampStream, err := breakpoint.NewBreakpointStream(ampPoints, wfmt.SampleRate)

	freqs, err := ioutil.ReadFile(*freqpoints)
	if err != nil {
		panic(err)
	}
	freqPoints, err := breakpoint.ParseBreakpoints(bytes.NewReader(freqs))
	if err != nil {
		panic(err)
	}
	freqStream, err := breakpoint.NewBreakpointStream(freqPoints, wfmt.SampleRate)
	if err != nil {
		panic(err)
	}
	// create wave file sampled at 44.1Khz w/ 16-bit frames

	frames := generate(*duration, stringToShape[*shape], ampStream, freqStream, wfmt)
	wave.WriteFrames(frames, wfmt, *output)
	fmt.Println("done")
}

func generate(dur int, shape synth.Shape, ampStream, freqStream *breakpoint.BreakpointStream, wfmt wave.WaveFmt) []wave.Frame {
	reqFrames := dur * wfmt.SampleRate
	frames := make([]wave.Frame, reqFrames)
	osc, err := synth.NewOscillator(wfmt.SampleRate, shape)
	if err != nil {
		panic(err)
	}

	for i := range frames {
		amp := ampStream.Tick()
		freq := freqStream.Tick()
		frames[i] = wave.Frame(amp * osc.Tick(freq))
	}

	return frames
}
