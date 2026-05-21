package ui

import "testing"

func TestNewLayoutFillsAvailableTerminalSpace(t *testing.T) {
	layout := NewLayout(120, 40)

	if !layout.Playable {
		t.Fatal("layout is not playable")
	}

	if layout.BoardWidth != 43 {
		t.Fatalf("board width = %d, want 43", layout.BoardWidth)
	}

	if layout.BoardHeight != 34 {
		t.Fatalf("board height = %d, want 34", layout.BoardHeight)
	}
}

func TestNewLayoutRejectsSmallTerminal(t *testing.T) {
	layout := NewLayout(40, 12)

	if layout.Playable {
		t.Fatal("layout is playable, want too small")
	}
}
