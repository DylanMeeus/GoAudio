package synthesizer

import (
	"math"
	"strings"
)

var (
	// noteIndex for use in calculations where a user passes a note
	noteIndex = map[string]int{
		"a":  0,
		"a#": 1,
		"bb": 1,
		"b":  2,
		"c":  3,
		"c#": 4,
		"db": 4,
		"d":  5,
		"d#": 6,
		"eb": 6,
		"e":  7,
		"f":  8,
		"f#": 9,
		"gb": 9,
		"g":  10,
		"g#": 11,
		"ab": 11,
	}
)

// ADSR creates an attack -> decay -> sustain -> release envelope
// time durations are passes as seconds.
// returns the value + the current time
func ADSR(maxamp, duration, attacktime, decaytime, sus, releasetime, controlrate float64, currentframe int) float64 {
	dur := duration * controlrate
	at := attacktime * controlrate
	dt := decaytime * controlrate
	rt := releasetime * controlrate
	cnt := float64(currentframe)

	amp := 0.0
	if cnt < dur {
		if cnt <= at {
			// attack
			amp = cnt * (maxamp / at)
		} else if cnt <= (at + dt) {
			// decay
			amp = ((sus-maxamp)/dt)*(cnt-at) + maxamp
		} else if cnt <= dur-rt {
			// sustain
			amp = sus
		} else if cnt > (dur - rt) {
			// release
			amp = -(sus/rt)*(cnt-(dur-rt)) + sus
		}
	}

	return amp
}

// NoteToFrequency turns a given note & octave into a frequency
// using Equal-Tempered tuning with reference pitch = A440
func NoteToFrequency(note string, octave int) float64 {
	// TODO: Allow for tuning systems other than Equal-Tempered A440?
	// clean the input
	note = strings.ToLower(strings.TrimSpace(note))
	ni := noteIndex[note]
	if ni >= 3 {
		// correct for octaves starting at C, not A.
		octave--
	}
	FR := 440.
	// we adjust the octave (-4) as the reference frequency is in the fourth octave
	// this effectively allows us to generate any octave above or below the reference octave
	return FR * math.Pow(2, float64(octave-4)+(float64(ni)/12.))
}
