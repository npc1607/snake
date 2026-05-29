package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"snake/internal/game"
	"snake/internal/session"
)

const frameDuration = 95 * time.Millisecond

type tickMsg time.Time

type stateMsg game.State

type clientConnectedMsg struct {
	client *session.Client
	err    error
}

func tick() tea.Cmd {
	return tea.Tick(frameDuration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func pollState(next func() (game.State, bool)) tea.Cmd {
	return func() tea.Msg {
		state, ok := next()
		if !ok {
			return nil
		}

		return stateMsg(state)
	}
}

func connectClient(address string) tea.Cmd {
	return func() tea.Msg {
		client, err := session.NewClient(address)
		return clientConnectedMsg{
			client: client,
			err:    err,
		}
	}
}
