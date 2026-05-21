package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"snake/internal/game"
)

type keyAction int

const (
	keyActionNone keyAction = iota
	keyActionQuit
	keyActionRestart
	keyActionPause
	keyActionTurnUp
	keyActionTurnDown
	keyActionTurnLeft
	keyActionTurnRight
)

func keyActionFor(msg tea.KeyPressMsg) keyAction {
	key := msg.Key()

	if key.Mod == tea.ModCtrl && key.Code == 'c' {
		return keyActionQuit
	}

	switch key.Code {
	case tea.KeyEsc:
		return keyActionQuit
	case tea.KeyUp:
		return keyActionTurnUp
	case tea.KeyDown:
		return keyActionTurnDown
	case tea.KeyLeft:
		return keyActionTurnLeft
	case tea.KeyRight:
		return keyActionTurnRight
	case tea.KeySpace:
		return keyActionPause
	}

	switch strings.ToLower(key.Text) {
	case "q":
		return keyActionQuit
	case "r":
		return keyActionRestart
	case " ", "p":
		return keyActionPause
	case "w", "k":
		return keyActionTurnUp
	case "s", "j":
		return keyActionTurnDown
	case "a", "h":
		return keyActionTurnLeft
	case "d", "l":
		return keyActionTurnRight
	default:
		return keyActionNone
	}
}

func directionForAction(action keyAction) (game.Direction, bool) {
	switch action {
	case keyActionTurnUp:
		return game.DirectionUp, true
	case keyActionTurnDown:
		return game.DirectionDown, true
	case keyActionTurnLeft:
		return game.DirectionLeft, true
	case keyActionTurnRight:
		return game.DirectionRight, true
	default:
		return game.DirectionRight, false
	}
}
