package tui

import (
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"

	"snake/internal/game"
	"snake/internal/ui"
)

// App adapts the game engine to Bubble Tea's Model interface.
type App struct {
	engine         *game.Engine
	renderer       ui.Renderer
	terminalWidth  int
	terminalHeight int
}

// NewApp creates a Bubble Tea model for the snake game.
func NewApp() App {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return App{
		engine:   game.NewEngine(game.DefaultConfig(), random),
		renderer: ui.NewRenderer(),
	}
}

// Init starts the recurring game tick.
func (a App) Init() tea.Cmd {
	return tick()
}

// Update handles terminal messages and advances the game.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.terminalWidth = msg.Width
		a.terminalHeight = msg.Height
		layout := ui.NewLayout(msg.Width, msg.Height)
		if layout.Playable {
			a.engine.Resize(layout.BoardWidth, layout.BoardHeight)
		}
	case tea.KeyPressMsg:
		action := keyActionFor(msg)
		switch action {
		case keyActionQuit:
			return a, tea.Quit
		case keyActionRestart:
			a.engine.Reset()
		case keyActionPause:
			a.engine.TogglePause()
		default:
			if direction, ok := directionForAction(action); ok {
				a.engine.Turn(direction)
			}
		}
	case tickMsg:
		a.engine.Tick()
		return a, tick()
	}

	return a, nil
}

// View renders the current game screen.
func (a App) View() tea.View {
	view := tea.NewView(a.renderer.Render(a.engine.State(), a.terminalWidth, a.terminalHeight))
	view.AltScreen = true
	view.WindowTitle = "Snake"

	return view
}
