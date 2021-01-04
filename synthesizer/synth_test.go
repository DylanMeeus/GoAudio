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

	parseNoteFrequencyTests = []struct {
		in  string
		out float64
	}{
		{
			"g2",
			98.,
		},
		{
			"a4",
			440.,
		},
		{
			"a7",
			3520.,
		},
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

// TestParseNoteFrequency tests a selection of strings containing note+octave and ensures that the
// output is the expected frequency
func TestParseNoteFrequency(t *testing.T) {
	for _, test := range parseNoteFrequencyTests {
		t.Run("", func(t *testing.T) {
			freq, err := synth.ParseNoteToFrequency(test.in)
			if err != nil {
				t.Fatalf("Unexpected error occurred: %v", err)
			}
			if !floatFuzzyEquals(freq, test.out) {
				t.Fatalf("Expected %v but got %v for (%s)", test.out, freq, test.in)
			}
		})

	}
}

// BenchmarkNoteToFrequency benchmarks the lookup of a frequency given a note & string
func BenchmarkNoteToFrequency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		synth.NoteToFrequency("A", 4)
	}
}

// BenchmarkParseNoteFrequency benchmarks the 'translation' of a note/octave given as a string
// to a frequency
func BenchmarkParseNoteFrequency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		synth.ParseNoteToFrequency("A4")
	}
}

// are the floats equal, within some grace region?
// to deal with floating point representation errors
func floatFuzzyEquals(f1, f2 float64) bool {
	return f1 > f2-10e-2 && f1 < f2+10e-2
}
