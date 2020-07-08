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
	internal.ReadFloatFrames(filename)
	fmt.Println("done")
}
