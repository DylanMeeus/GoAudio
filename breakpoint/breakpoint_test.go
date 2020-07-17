package util

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
)

func TestBreakpoint(t *testing.T) {
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

func TestValueAt(t *testing.T) {
	for _, test := range valueAtTests {
		t.Run("", func(t *testing.T) {
			res := ValueAt(test.breaks, test.time)
			if res != test.out {
				t.Fatalf("expected %v, got %v", test.out, res)
			}
		})
	}

}
