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
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Should not have panic'd!\n%v", r)
		}
	}()
	goldenfile := "./golden/maybe-next-time.wav"
	wav, err := ReadWaveFile(goldenfile)
	if err != nil {
		t.Fatalf("Should be able to read wave file: %v", err)
	}

	// assert that the wav file looks as expected.
	if wav.SampleRate != 44100 {
		t.Fatalf("Expected SR 44100, got: %v", wav.SampleRate)
	}

	if wav.NumChannels != 2 {
		t.Fatalf("Expected 2 channels, got: %v", wav.NumChannels)
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
