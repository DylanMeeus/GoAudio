package wave

import (
	"testing"
)

var (
	testBatchSamples = []struct {
		wave     Wave
		timespan float64
		out      [][]Frame
	}{
		{
			Wave{},
			0,
			[][]Frame{[]Frame{}},
		},
		{
			Wave{
				WaveFmt: WaveFmt{
					SampleRate: 2, // 2 seconds per sample
				},
				WaveData: WaveData{
					Frames: makeSampleSlice(1, 2, 3, 4, 5, 6, 7, 8),
				},
			},
			2,
			[][]Frame{makeSampleSlice(1, 2, 3, 4), makeSampleSlice(5, 6, 7, 8)},
		},
		{
			Wave{
				WaveFmt: WaveFmt{
					SampleRate: 2, // 2 seconds per sample
				},
				WaveData: WaveData{
					Frames: makeSampleSlice(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
				},
			},
			2,
			[][]Frame{makeSampleSlice(1, 2, 3, 4), makeSampleSlice(5, 6, 7, 8), makeSampleSlice(9, 10)},
		},
	}
)

func TestBatching(t *testing.T) {
	t.Logf("Testing batching of samples per time slice")
	for _, test := range testBatchSamples {
		t.Run("", func(t *testing.T) {
			res := BatchSamples(test.wave, test.timespan)
			if !compareSampleSlices(res, test.out) {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}
}

// helper functions for testing

func makeSampleSlice(input ...float64) (out []Frame) {
	for _, f := range input {
		out = append(out, Frame(f))
	}
	return
}

// compareSampleSlices makes sure both slices are the same
func compareSampleSlices(a, b [][]Frame) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		for j, v := range x {
			if b[i][j] != v {
				return false
			}
		}
	}
	return true
}
