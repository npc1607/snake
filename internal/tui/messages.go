package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

const frameDuration = 95 * time.Millisecond

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(frameDuration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
