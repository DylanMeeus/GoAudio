package wave

import (
	"testing"
)

var (
	scaleFrameTests = []struct {
		input int
		bits  int
		out   Frame
	}{
		{
			0,
			16,
			0,
		},
		{
			127,
			8,
			1,
		},
		{
			-127,
			8,
			-1,
		},
		{
			32_767,
			16,
			1,
		},
		{
			32_767,
			16,
			1,
		},
		{
			0,
			16,
			0,
		},
		{
			-32_767,
			16,
			-1,
		},
	}

	// readFileTest to iterate over various files and make sure the output meets
	// the expected Wave struct.
	readFileTests = []struct {
		file     string
		expected Wave
	}{}
)

// TestReadFile reads a golden file and ensures that is parsed as expected
func TestReadFile(t *testing.T) {
	tests := []struct {
		fileName   string
		sampleRate int
		channels   int
	}{
		{
			fileName:   "./golden/maybe-next-time.wav",
			sampleRate: 44100,
			channels:   2,
		},
		{
			fileName:   "./golden/c2-24bit.wav",
			sampleRate: 44100,
			channels:   2,
		},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Should not have panic'd!\n%v", r)
		}
	}()

	for _, tt := range tests {
		wav, err := ReadWaveFile(tt.fileName)
		if err != nil {
			t.Fatalf("Should be able to read wave file: %v", err)
		}
		if wav.SampleRate != tt.sampleRate {
			t.Fatalf("Expected SR %d, got: %v", tt.sampleRate, wav.SampleRate)
		}
		if wav.NumChannels != tt.channels {
			t.Fatalf("Expected %d channels, got: %v", tt.channels, wav.NumChannels)
		}
	}
}

func TestScaleFrames(t *testing.T) {
	for _, test := range scaleFrameTests {
		t.Run("", func(t *testing.T) {
			res := scaleFrame(test.input, test.bits)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}
}
