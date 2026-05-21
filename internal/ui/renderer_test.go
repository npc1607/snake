package ui

import (
	"testing"

	"charm.land/lipgloss/v2"

	"snake/internal/game"
)

func TestRendererFitsTerminal(t *testing.T) {
	renderer := NewRenderer()
	state := game.NewEngine(game.Config{Width: 43, Height: 34}, nil).State()

	output := renderer.Render(state, 120, 40)
	width, height := lipgloss.Size(output)

	if width > 120 {
		t.Fatalf("rendered width = %d, want <= 120", width)
	}

	if height > 40 {
		t.Fatalf("rendered height = %d, want <= 40", height)
	}
}
