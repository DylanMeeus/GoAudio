package util

import (
	"strings"
	"testing"
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
