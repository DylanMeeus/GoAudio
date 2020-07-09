package main

import (
	"fmt"
	"github.com/DylanMeeus/Audio/wave/internal"
	"os"
)

func main() {
	fmt.Println("Parsing wave file..")
	args := os.Args
	if len(args) < 1 {
		panic("please provide a file..")
	}
	filename := args[1]
	wave, err := internal.ReadWaveFile(filename)
	if err != nil {
		panic("Could not parse wave file")
	}
	fmt.Printf("Wave samples: %v\n", wave.Samples)
	fmt.Println("done")
}
