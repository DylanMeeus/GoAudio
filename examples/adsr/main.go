package main

import (
	"flag"
	"fmt"

	synth "github.com/DylanMeeus/GoAudio/synthesizer"
	"github.com/DylanMeeus/GoAudio/wave"
)

func main() {
	flag.Parse()

	osc, err := synth.NewOscillator(44100, synth.SINE)
	if err != nil {
		panic(err)
	}

	sr := 44100
	duration := sr * 10

	frames := []wave.Frame{}
	var adsrtime int
	for i := 0; i < duration; i++ {
		value, time := synth.ADSR(1, 10, 1, 1, 0.7, 5, float64(sr), adsrtime)
		adsrtime = time
		frames = append(frames, wave.Frame(value*osc.Tick(440)))
	}

	wfmt := wave.NewWaveFmt(1, 1, sr, 16, nil)
	wave.WriteFrames(frames, wfmt, "output.wav")

	fmt.Println("done writing to output.wav")
}
