package main

// tool to print the interpreted the content of a .wave file

import (
	"flag"
	"fmt"
	"os"

	"github.com/DylanMeeus/GoAudio/wave"
	wav "github.com/DylanMeeus/GoAudio/wave"
)

var (
	input       = flag.String("i", "", "input file")
	withSamples = flag.Bool("s", false, "with raw audio samples")
)

// printHeader
func printHeader(h wave.WaveHeader) {
	fmt.Println("Header")
	fmt.Printf("Chunk ID: %v\n", string(h.ChunkID))
	fmt.Printf("Chunk Size: %v\n", h.ChunkSize)
	fmt.Printf("Format: %v\n", string(h.Format))

}

func printFormat(f wave.WaveFmt) {
	fmt.Println("Wave FMT")
	fmt.Printf("SubchunkID: %v\n", string(f.Subchunk1ID))
	fmt.Printf("SubchunkSize: %v\n", f.Subchunk1Size)
	fmt.Printf("AudioFormat: %v\n", f.AudioFormat)
	fmt.Printf("Channels: %v\n", f.NumChannels)
	fmt.Printf("SampleRate: %v\n", f.SampleRate)
	fmt.Printf("ByteRate: %v\n", f.ByteRate)
	fmt.Printf("BlockAlign: %v\n", f.BlockAlign)
	fmt.Printf("BitsPerSample: %v\n", f.BitsPerSample)

}

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

	printHeader(wave.WaveHeader)
	fmt.Printf("===============\n")
	printFormat(wave.WaveFmt)
	if ws {
		fmt.Printf("===============\nSamples:\n")
		fmt.Printf("%v\n", wave.RawData)
	}

}
