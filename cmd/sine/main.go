package main

/* program to create a pitch perfect (440Hz) sound */

import (
	"fmt"
	"math"
)

const (
	nsamps = 50
)

var (
	tau = math.Pi * 2
)

func main() {
	fmt.Printf("generating sine wave..\n")
	generate()
	fmt.Printf("done")
}

func generate() {
	var angleincr float64 = tau / nsamps
	for i := 0; i < nsamps; i++ {
		samp := math.Sin(angleincr * float64(i))
		fmt.Printf("%.81f\t%81.f\n", samp, samp)
	}
}
