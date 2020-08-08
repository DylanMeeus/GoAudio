package main

import (
	"flag"
	"io/ioutil"
	"strconv"
	"strings"

	wav "github.com/DylanMeeus/GoAudio/wave"
)

// program to extract breakpoint data from an input source.

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	window = flag.Int("w", 15, "window of time for capturing breakpoint data")
)

func main() {
	flag.Parse()
	infile := *input
	outfile := *output
	wave, err := wav.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	ticks := float64(*window) / 1000.0
	batches := wav.BatchSamples(wave, ticks)

	strout := strings.Builder{}
	elapsed := 0.0
	for _, b := range batches {
		maxa := maxAmp(b)
		es := strconv.FormatFloat(elapsed, 'f', 8, 64)
		fs := strconv.FormatFloat(maxa, 'f', 8, 64)
		strout.WriteString(es + ":" + fs + "\n")
		elapsed += ticks
	}

	err = ioutil.WriteFile(outfile, []byte(strout.String()), 0644)
	if err != nil {
		panic(err)
	}

}

// return the maximum amplitude of these samples
func maxAmp(ss []wav.Frame) float64 {
	if len(ss) == 0 {
		return 0
	}
	max := -1.0 // because they are in range -1 .. 1
	for _, a := range ss {
		if float64(a) > max {
			max = float64(a)
		}
	}
	return float64(max)
}
