package main

import (
	"flag"
	"fmt"
	"github.com/DylanMeeus/GoAudio/wave/pkg"
)

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	pan    = flag.Float64("p", 0.0, "pan in range of -1 (left) to 1 (right)")
)

type panposition struct {
	left, right float64
}

func calculatePosition(position float64) panposition {
	position *= 0.5
	return panposition{
		left:  position - 0.5,
		right: position + 0.5,
	}
}

func main() {
	fmt.Println("Parsing wave file..")
	flag.Parse()
	infile := *input
	outfile := *output
	panfac := *pan
	wave, err := pkg.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	fmt.Printf("Read %v samples\n", len(wave.Samples))

	// now try to write this file
	pos := calculatePosition(panfac)
	fmt.Printf("panfac: %v\npanpos: %v\n", panfac, pos)
	scaledSamples := applyPan(wave.Samples, calculatePosition(panfac))
	wave.NumChannels = 2 // samples are now stereo, so we need dual channels
	if err := pkg.WriteSamples(scaledSamples, wave.WaveFmt, outfile); err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func applyPan(samples []pkg.Sample, p panposition) []pkg.Sample {
	out := []pkg.Sample{}
	for _, s := range samples {
		out = append(out, pkg.Sample(float64(s)*p.left))
		out = append(out, pkg.Sample(float64(s)*p.right))
	}
	return out
}
