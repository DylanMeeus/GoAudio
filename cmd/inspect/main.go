package main

// tool to print the interpreted the content of a .wave file

import (
	"flag"
	"fmt"
	wav "github.com/DylanMeeus/GoAudio/wave"
	"os"
)

var (
	input       = flag.String("i", "", "input file")
	withSamples = flag.Bool("s", false, "with raw audio samples")
)

func main() {
	flag.Parse()
	as := os.Args[1:]
	if len(as) == 0 {
		panic("Please provide file")
	}
	infile := as[0]
	ws := *withSamples
	wave, err := wav.ReadWaveFile(infile)
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
