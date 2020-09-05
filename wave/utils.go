package wave

// utility functions for dealing with wave files.

// BatchSamples batches the samples per requested timespan expressed in seconds
func BatchSamples(data Wave, seconds float64) [][]Frame {
	if seconds == 0 {
		return [][]Frame{
			data.Frames,
		}
	}

	samples := data.Frames

	sampleSize := int(float64(data.SampleRate*data.NumChannels) * float64(seconds))

	batches := len(samples) / sampleSize
	if len(samples)%sampleSize != 0 {
		batches++
	}

	batched := make([][]Frame, batches) // this should be round up..
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

// FloatsToFrames turns a slice of float64 to a slice of frames
func FloatsToFrames(fs []float64) []Frame {
	frames := make([]Frame, len(fs))
	for i, f := range fs {
		frames[i] = Frame(f)
	}
	return frames
}
