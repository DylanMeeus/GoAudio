package main

import (
	"flag"
	"fmt"
	brk "github.com/DylanMeeus/GoAudio/breakpoint"
	wav "github.com/DylanMeeus/GoAudio/wave"
	"math"
	"os"
)

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	pan    = flag.Float64("p", 0.0, "pan in range of -1 (left) to 1 (right)")
	brkpnt = flag.String("b", "", "breakpoint file")
)

type panposition struct {
	left, right float64
}

// calculateConstantPowerPosition finds the position of each speaker using a constant power function
func calculateConstantPowerPosition(position float64) panposition {
	// half a sinusoid cycle
	var halfpi float64 = math.Pi / 2
	r := math.Sqrt(2.0) / 2

	// scale position to fit in this range
	scaled := position * halfpi

	// each channel uses 1/4 of a cycle
	angle := scaled / 2
	pos := panposition{}
	pos.left = r * (math.Cos(angle) - math.Sin(angle))
	pos.right = r * (math.Cos(angle) + math.Sin(angle))
	return pos
}

func calculatePosition(position float64) panposition {
	position *= 0.5
	return panposition{
		left:  position - 0.5,
		right: position + 0.5,
	}
}

// set a single pan (without breakpoints)
func setPan() {
	fmt.Println("Parsing wave file..")
	flag.Parse()
	infile := *input
	outfile := *output
	panfac := *pan
	wave, err := wav.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	fmt.Printf("Read %v samples\n", len(wave.Frames))

	// now try to write this file
	pos := calculatePosition(panfac)
	fmt.Printf("panfac: %v\npanpos: %v\n", panfac, pos)
	scaledFrames := applyPan(wave.Frames, calculatePosition(panfac))
	wave.NumChannels = 2 // samples are now stereo, so we need dual channels
	if err := wav.WriteFrames(scaledFrames, wave.WaveFmt, outfile); err != nil {
		panic(err)
	}
	fmt.Println("done")
}

func withBreakpointFile() {
	// apply a stereo pan by applying a breakpoint file using linear interpolation
	flag.Parse()

	file, err := os.Open(*brkpnt)
	if err != nil {
		panic(err)
	}
	pnts, err := brk.ParseBreakpoints(file)
	if err != nil {
		panic(err)
	}

	// verify that the breakpoints are valid
	if brk.Breakpoints(pnts).Any(func(b brk.Breakpoint) bool {
		return b.Value > 1 || b.Value < -1
	}) {
		panic("Breakpoint file contains invalid values")
	}

	// read the content file
	infile := *input
	wave, err := wav.ReadWaveFile(infile)

	// now we want to apply the pan per sample
	// for this we need to know the time at each sample.
	// which is the reciprocal of sample rate
	timeincr := 1.0 / float64(wave.SampleRate)
	var frametime float64
	inframes := wave.Frames
	var out []wav.Frame

	wave.WaveFmt.SetChannels(2)
	for _, s := range inframes {
		// apply pan
		_, pos := brk.ValueAt(pnts, frametime, 0)
		pan := calculateConstantPowerPosition(pos)
		out = append(out, wav.Frame(float64(s)*pan.left))
		out = append(out, wav.Frame(float64(s)*pan.right))
		frametime += timeincr
	}

	wav.WriteFrames(out, wave.WaveFmt, *output)
}

func main() {
	//setPan()
	withBreakpointFile()
}

func applyPan(samples []wav.Frame, p panposition) []wav.Frame {
	out := []wav.Frame{}
	for _, s := range samples {
		out = append(out, wav.Frame(float64(s)*p.left))
		out = append(out, wav.Frame(float64(s)*p.right))
	}
	return out
}
