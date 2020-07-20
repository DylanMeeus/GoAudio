package wave

import (
	"time"
)

// utility functions for dealing with wave files.

// BatchSamples batches the samples per requested timespan (timespan in seconds)
func BatchSamples(data Wave, timespan time.Duration) [][]Sample {
	if timespan == 0 {
		return [][]Sample{
			data.Samples,
		}
	}

	samples := data.Samples

	sampleSize := data.SampleRate * int(timespan)
	// now we need to grab these slices..

	batched := make([][]Sample, len(samples)/sampleSize)
	for i := 0; i < len(batched); i++ {
		start := i * sampleSize
		if start > len(samples) {
			return batched
		}
		maxTake := i*sampleSize + sampleSize
		if maxTake >= len(samples)-1 {
			maxTake = len(samples)
		}
		subs := samples[start:maxTake]
		batched[i] = subs
	}
	// figure out how many samples per duration?
	// depends on the samplerate, which is 'samples per second'
	return batched
}
