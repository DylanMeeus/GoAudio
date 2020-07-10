package main

import (
	"flag"
	"fmt"
	"github.com/DylanMeeus/Audio/wave/internal"
)

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	amp    = flag.Float64("a", 1.0, "amp mod factor")
)

func main() {
	fmt.Println("Parsing wave file..")
	flag.Parse()
	infile := *input
	outfile := *output
	scale := *amp
	wave, err := internal.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	fmt.Printf("Read %v samples\n", len(wave.Samples))

	// now try to write this file
	scaledSamples := changeAmplitude(wave.Samples, scale)
	if err := internal.WriteSamples(scaledSamples, wave.WaveFmt, outfile); err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func changeAmplitude(samples []internal.Sample, scalefactor float64) []internal.Sample {
	for i, s := range samples {
		samples[i] = internal.Sample(float64(s) * scalefactor)
	}
	return samples
}
