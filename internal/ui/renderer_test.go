package ui

import (
	"strings"
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

func TestRendererShowsRoomMembersInSidebar(t *testing.T) {
	renderer := NewRenderer()
	state := game.NewEngine(game.Config{Width: 43, Height: 34, PlayerCount: 2}, nil).State()

	output := renderer.RenderRoom(state, 120, 40, []RoomMember{
		{Name: "Player 1", Addr: "[::]:7777", Status: "hosting", Current: true},
		{Name: "Player 2", Addr: "127.0.0.1:7777", Status: "joined", LatencyMS: 42},
	})

	if !strings.Contains(output, "ROOM") {
		t.Fatal("rendered output does not include room section")
	}
	if !strings.Contains(output, "Player 2") || !strings.Contains(output, "joined") {
		t.Fatal("rendered output does not include joined player")
	}
	if !strings.Contains(output, "Latency 42ms") {
		t.Fatal("rendered output does not include player latency")
	}
}
