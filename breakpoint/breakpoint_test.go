package breakpoint

import (
	"strings"
	"testing"
)

var (
	valueAtTests = []struct {
		breaks []Breakpoint
		time   float64
		out    float64
	}{
		{
			[]Breakpoint{},
			0,
			0,
		},
		{
			[]Breakpoint{},
			0,
			0,
		},
		{
			[]Breakpoint{
				{
					Time:  0,
					Value: 1,
				},
			},
			0,
			1,
		},
		{
			[]Breakpoint{
				{
					Time:  0,
					Value: 1,
				},
			},
			2,
			1,
		},
		{
			[]Breakpoint{
				{
					Time:  0,
					Value: 2,
				},
				{
					Time:  1,
					Value: 0,
				},
			},
			2,
			0,
		},
		{
			[]Breakpoint{
				{
					Time:  1,
					Value: 0,
				},
				{
					Time:  2,
					Value: 10,
				},
				{
					Time:  3,
					Value: 100,
				},
			},
			1.5, // linear interpolation should give this result
			5,
		},
	}

	anyTests = []struct {
		in   Breakpoints
		anyf func(Breakpoint) bool
		out  bool
	}{
		{
			in:   Breakpoints{Breakpoint{1, 0}, Breakpoint{2, 5}, Breakpoint{3, 10}},
			anyf: func(b Breakpoint) bool { return b.Value > 5 },
			out:  true,
		},
		{
			in:   Breakpoints{Breakpoint{1, 0}, Breakpoint{2, 5}, Breakpoint{3, 10}},
			anyf: func(b Breakpoint) bool { return b.Value > 15 },
			out:  false,
		},
	}
)

func TestBreakpoint(t *testing.T) {
	// TODO: these breakpoints are invalid as long as we can't timetravel. Time should be strictly
	// increasing..
	input := `
	3.0:1.31415
	-1:1.32134
	0:0
	`
	brk, err := ParseBreakpoints(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Should be able to parse breakpoints: %v", err)
	}

	if brk[0].Time != 3.0 {
		t.Fatalf("Breakpoint does not match input, expected %v, got %v", 3.0, brk[0].Time)
	}
	if brk[1].Time != -1 {
		t.Fatal("Breakpoint does not match input")
	}
	if brk[2].Time != 0 {
		t.Fatal("Breakpoint does not match input")
	}

	if brk[0].Value != 1.31415 {
		t.Fatal("Breakpoint does not match input")
	}
	if brk[1].Value != 1.32134 {
		t.Fatal("Breakpoint does not match input")
	}
	if brk[2].Value != 0 {
		t.Fatal("Breakpoint does not match input")
	}
}

// TestbreakpointStream tests various assumptions of how ticks should be handled
// We test at 10 samples per second for convenience (easier to verify correctness)
// Thus each -tick- should update our position by 0.1
func TestBreakpointStream(t *testing.T) {
	input := `
	0:100
	10:50
	20:100
	`
	samplerate := 10 // 10 sample per second
	brks, err := ParseBreakpoints(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Should be able to parse breakpoints: %v", err)
	}
	stream, err := NewBreakpointStream(brks, samplerate)
	if err != nil {
		t.Fatalf("Should be able to create breakpoint stream: %v", err)
	}

	if stream.Increment != 1.0/float64(samplerate) {
		t.Fatal("Incorrect increment")
	}
	value := stream.Tick()
	if stream.CurrentPosition != 0.1 {
		t.Fatal("Incorrectly incremented stream")
	}
	// move to 5 seconds (50 ticks)
	for i := 0; i < 50; i++ {
		value = stream.Tick()
	}
	if value != 75 {
		t.Fatalf("Expected 75 but got %v", value)
	}

	// move to second 15
	for i := 0; i < 100; i++ {
		value = stream.Tick()
	}

	t.Log("Should be able to get values beyond the end of the stream")
	for i := 0; i < 10e5; i++ {
		value = stream.Tick()
	}
	if value != 100 {
		t.Fatalf("Value should be 100, but got %v", value)
	}
}

func TestValueAt(t *testing.T) {
	for _, test := range valueAtTests {
		t.Run("", func(t *testing.T) {
			_, res := ValueAt(test.breaks, test.time, 0)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}
}

func TestAny(t *testing.T) {
	for _, test := range anyTests {
		t.Run("", func(t *testing.T) {
			out := test.in.Any(test.anyf)
			if out != test.out {
				t.Fatalf("Expected %v but got %v", test.out, out)
			}
		})
	}

}
