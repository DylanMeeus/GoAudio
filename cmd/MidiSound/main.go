package main

/*
   Play sound based on midi notes..
*/

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
)

type config struct {
	Duration   int
	MidiNote   int
	SampleRate int
	Amplitude  float64
}

// Constants for generating our output
const (
	DURATION = iota
	MIDI
	SR
	AMP
	OUTFILE
	TOTAL_ARGS
)

const concertA = float64(440)

// usage: go run main.go [dur] [hz] [sr] [amp]
// e.g: go run main.go 2 440 44100 1
func main() {
	fmt.Printf("Generating..\n")
	generate(parseInput())
	fmt.Printf("Done\n")
}

// midi2hertz turns a midi note into a frequency
func midi2hertz(n int) float64 {
	middleC := concertA // oft used default for middleC
	return math.Pow(2, (float64(n)-69)/12) * middleC
}

// generate the sample based on the config
func generate(c config) {
	var (
		start float64 = 1.0
		end   float64 = 1.0e-4
		tau           = math.Pi * 2
	)

	hertz := midi2hertz(c.MidiNote)
	fmt.Printf("don't hertz me: %v\n", hertz)

	// setup output file
	file := os.Args[1:][OUTFILE]
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nsamples := c.Duration * c.SampleRate
	angleincr := tau * hertz / float64(nsamples)
	decayfac := math.Pow(end/start, 1.0/float64(nsamples))

	for i := 0; i < nsamples; i++ {
		sample := c.Amplitude * math.Sin(angleincr*float64(i))
		sample *= start
		start *= decayfac
		var buf [8]byte
		// I know my system is LittleEndian, use BigEndian if yours is not..
		binary.LittleEndian.PutUint32(buf[:], math.Float32bits(float32(sample)))
		bw, err := f.Write(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("\rWrote: %v bytes to %s", bw, file)
	}
	fmt.Printf("\n")
}

// will ignore error handling for now.. ;-)
func parseInput() config {
	args := os.Args[1:]
	if len(args) != TOTAL_ARGS {
		return config{}
	}
	dur, _ := strconv.Atoi(args[DURATION])
	midi, _ := strconv.Atoi(args[MIDI])
	sr, _ := strconv.Atoi(args[SR])
	amp, _ := strconv.ParseFloat(args[AMP], 64)
	return config{
		Duration:   dur,
		MidiNote:   midi,
		SampleRate: sr,
		Amplitude:  amp,
	}
}
