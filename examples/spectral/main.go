package main

import (
	"fmt"

	audiomath "github.com/DylanMeeus/GoAudio/math"
	"github.com/DylanMeeus/GoAudio/wave"
	"math"
)

var _ = audiomath.HFFT

func main() {
	w, err := wave.ReadWaveFile("frerejacques.wav")
	if err != nil {
		panic(err)
	}
	fmt.Println(2 * math.Pi)
	c128 := audiomath.FFT(w.Frames)
	fmt.Printf("%v\n", c128)
}
