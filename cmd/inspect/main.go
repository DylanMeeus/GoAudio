package main

// tool to print the interpreted the content of a .wave file

import (
	"flag"
	"fmt"
	"github.com/DylanMeeus/Audio/wave/internal"
)

var (
	input       = flag.String("i", "", "input file")
	withSamples = flag.Bool("s", false, "with raw audio samples")
)

func main() {
	flag.Parse()
	infile := *input
	ws := *withSamples
	wave, err := internal.ReadWaveFile(infile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Header:\n")
	fmt.Printf("%+v\n", wave.WaveHeader)
	fmt.Printf("===============\n")
	fmt.Printf("Fmt:\n")
	fmt.Printf("%+v\n", wave.WaveFmt)
	if ws {
		fmt.Printf("===============\nSamples:\n")
		fmt.Printf("%v\n", wave.RawData)
	}

}
