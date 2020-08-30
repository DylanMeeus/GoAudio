package wave

import (
	"testing"
)

var (
	rescaleFrameTests = []struct {
		input Frame
		bits  int
		out   int
	}{
		{
			Frame(1),
			8,
			127,
		},
		{
			Frame(-1),
			8,
			-127,
		},
		{
			Frame(1),
			16,
			32_767,
		},
		{
			Frame(-1),
			16,
			-32_767,
		},
	}
)

func TestRescaleFrames(t *testing.T) {
	for _, test := range rescaleFrameTests {
		t.Run("", func(t *testing.T) {
			res := rescaleFrame(test.input, test.bits)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}
}

// TestWriteWave reads wave file and writes it, ensuring nothing is different between the two
func TestWriteWave(t *testing.T) {
	goldenfile := "./golden/maybe-next-time.wav"
	wav, err := ReadWaveFile(goldenfile)
	if err != nil {
		t.Fatalf("Should be able to read wave file: %v", err)
	}

	if err := WriteFrames(wav.Frames, wav.WaveFmt, "output.wav"); err != nil {
		t.Fatalf("Should be able to write file: %v", err)
	}
}
