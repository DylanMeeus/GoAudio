package breakpoint

// convenience functions for dealing with breakpoints

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// Breakpoints is a collection of Breakpoint (time:value) pairs
type Breakpoints []Breakpoint

// Breakpoint represents a time:value pair
type Breakpoint struct {
	Time, Value float64
}

// BreakpointStream can be used to to treat breakpoints as a stream of data
// Each 'tick' can manipulate the state of the breakpoint stream
type BreakpointStream struct {
	Breakpoints     Breakpoints
	Left            Breakpoint
	Right           Breakpoint
	IndexLeft       int
	IndexRight      int
	CurrentPosition float64 // current position in timeframes
	Increment       float64
	Width           float64
	Height          float64
	HasMore         bool
}

// Tick returns the next value in the breakpoint stream
func (b *BreakpointStream) Tick() (out float64) {
	if !b.HasMore {
		// permanently the last value
		return b.Right.Value
	}
	if b.Width == 0.0 {
		out = b.Right.Value
	} else {
		// figure out value from linear interpolation
		frac := (float64(b.CurrentPosition) - b.Left.Time) / b.Width
		out = b.Left.Value + (b.Height * frac)
	}

	// prepare for next frame
	b.CurrentPosition += b.Increment
	if b.CurrentPosition > b.Right.Time {
		// move to next span
		b.IndexLeft++
		b.IndexRight++
		if b.IndexRight < len(b.Breakpoints) {
			b.Left = b.Breakpoints[b.IndexLeft]
			b.Right = b.Breakpoints[b.IndexRight]
			b.Width = b.Right.Time - b.Left.Time
			b.Height = b.Right.Value - b.Left.Value
		} else {
			// no more points
			b.HasMore = false
		}
	}
	return out
}

// NewBreakpointStream represents a slice of breakpoints streamed at a given sample rate
func NewBreakpointStream(bs []Breakpoint, sr int) (*BreakpointStream, error) {
	if len(bs) == 0 {
		return nil, errors.New("Need at least two points to create a stream")
	}
	left, right := bs[0], bs[1]
	return &BreakpointStream{
		Breakpoints:     Breakpoints(bs),
		Increment:       1.0 / float64(sr),
		IndexLeft:       0,
		IndexRight:      1,
		CurrentPosition: 0,
		Left:            left,
		Right:           right,
		Width:           right.Time - left.Time,   // first span
		Height:          right.Value - left.Value, // diff of first span
		HasMore:         len(bs) > 0,
	}, nil
}

// ParseBreakpoints reads the breakpoints from an io.Reader
// and turns them into a slice.
// A file is expected to be [time: value] formatted
// Will panic if file format is wrong
// TODO: don't panic
func ParseBreakpoints(in io.Reader) ([]Breakpoint, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	brkpnts := []Breakpoint{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return brkpnts, err
		}
		time := parts[0]
		value := parts[1]

		tf, err := strconv.ParseFloat(time, 64)
		if err != nil {
			return brkpnts, err
		}
		vf, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return brkpnts, err
		}

		brkpnts = append(brkpnts, Breakpoint{
			Time:  tf,
			Value: vf,
		})

	}
	return brkpnts, nil
}

// ValueAt returns the expected value at a given time (expressed as float64) by linear interpolation
// Returns the index at which we found our value as well as the value itself.
func ValueAt(bs []Breakpoint, time float64, startIndex int) (index int, value float64) {
	if len(bs) == 0 {
		return 0, 0
	}
	npoints := len(bs)

	// first we need to find a span containing our timeslot
	startSpan := startIndex // start of span
	for _, b := range bs[startSpan:] {
		if b.Time > time {
			break
		}
		startSpan++
	}

	// Our span is never-ending (the last point in our breakpoint file was hit)
	if startSpan == npoints {
		return startSpan, bs[startSpan-1].Value
	}

	left := bs[startSpan-1]
	right := bs[startSpan]

	// check for instant jump
	// 2 points having the same time...
	width := right.Time - left.Time

	if width == 0 {
		return startSpan, right.Value
	}

	frac := (time - left.Time) / width

	val := left.Value + ((right.Value - left.Value) * frac)

	return startSpan, val
}

// MinMaxValue returns the smallest and largest value found in the breakpoint file
func MinMaxValue(bs []Breakpoint) (smallest float64, largest float64) {
	// TODO: implement as SORT and return the first and last element
	if len(bs) == 0 {
		return
	}
	smallest = bs[0].Value
	largest = bs[0].Value
	for _, b := range bs[1:] {
		if b.Value < smallest {
			smallest = b.Value
		} else if b.Value > largest {
			largest = b.Value
		} else {
			// no op
		}
	}
	return
}

// Any returns true if any breakpoint matches the filter.
func (bs Breakpoints) Any(f func(Breakpoint) bool) bool {
	for _, b := range bs {
		if f(b) {
			return true
		}
	}
	return false
}
