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
		{"bb", 5, 932.32},
		{"g", 2, 98.},
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
			if freq != test.output {
				t.Fatalf("Expected %v but got %v for (%s%v)", test.output, freq, test.note, test.octave)
			}
		})

	}
}
