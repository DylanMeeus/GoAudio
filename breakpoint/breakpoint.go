package util

// convenience functions for dealing with breakpoints

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Breakpoints []Breakpoint

type Breakpoint struct {
	Time, Value float64
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
func ValueAt(bs []Breakpoint, time float64) float64 {
	if len(bs) == 0 {
		return 0
	}
	npoints := len(bs)

	// first we need to find a span containing our timeslot
	startSpan := 0 // start of span
	for _, b := range bs {
		if b.Time > time {
			break
		}
		startSpan++
	}
	// Our span is never-ending (the last point in our breakpoint file was hit)
	if startSpan == npoints {
		return bs[startSpan-1].Value
	}

	left := bs[startSpan-1]
	right := bs[startSpan]

	// check for instant jump
	// 2 points having the same time...
	width := right.Time - left.Time

	if width == 0 {
		return right.Value
	}

	frac := (time - left.Time) / width

	val := left.Value + ((right.Value - left.Value) * frac)

	return val
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
