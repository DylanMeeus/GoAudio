package synthesizer_test

import (
	"testing"

	synth "github.com/DylanMeeus/GoAudio/synthesizer"
)

var (
	// we take a random sampling of notes -> frequencies to avoid dealing with odd Floating Point
	// values
	noteFrequencyTests = []struct {
		note   string
		octave int
		output float64
	}{
		{"g", 2, 98.},
		{"bb", 5, 932.32},
		{"G", 2, 98.},
		{" G ", 2, 98.},
		{"a", 7, 3520.},
		{"a", 4, 440.},
		{"a", 2, 110.},
	}
)

// TestNoteToFrequency tests a selection of notes + octaves and verifies their frequency
func TestNoteToFrequency(t *testing.T) {
	for _, test := range noteFrequencyTests {
		t.Run("", func(t *testing.T) {
			freq := synth.NoteToFrequency(test.note, test.octave)
			if !floatFuzzyEquals(freq, test.output) {
				t.Fatalf("Expected %v but got %v for (%s%v)", test.output, freq, test.note, test.octave)
			}
		})

	}
}

// are the floats equal, within some grace region?
// to deal with floating point representation errors
func floatFuzzyEquals(f1, f2 float64) bool {
	return f1 > f2-10e-2 && f1 < f2+10e-2
}
