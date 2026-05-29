package tui

import (
	"fmt"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"

	"snake/internal/game"
	"snake/internal/session"
	"snake/internal/ui"
)

// App adapts the game engine to Bubble Tea's Model interface.
type App struct {
	engine         *game.Engine
	renderer       ui.Renderer
	terminalWidth  int
	terminalHeight int
	mode           string
	joinAddr       string
	hostSession    *session.Host
	clientSession  *session.Client
	remoteState    game.State
	err            error
}

// NewApp creates a Bubble Tea model for the snake game.
func NewApp() App {
	return NewLocalApp()
}

// NewLocalApp creates a Bubble Tea model for single-player mode.
func NewLocalApp() App {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return App{
		engine:   game.NewEngine(game.DefaultConfig(), random),
		renderer: ui.NewRenderer(),
		mode:     "local",
	}
}

// NewHostApp creates a Bubble Tea model for host mode.
func NewHostApp(listenAddr string) App {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	config := game.DefaultConfig()

	hostSession, err := session.NewHost(listenAddr)

	return App{
		engine:      game.NewEngine(config, random),
		renderer:    ui.NewRenderer(),
		mode:        "host",
		hostSession: hostSession,
		err:         err,
	}
}

// NewJoinApp creates a Bubble Tea model for client mode.
func NewJoinApp(joinAddr string) App {
	return App{
		renderer: ui.NewRenderer(),
		mode:     "join",
		joinAddr: joinAddr,
	}
}

// Init starts the recurring game tick.
func (a App) Init() tea.Cmd {
	if a.mode == "join" {
		return connectClient(a.joinAddr)
	}

	return tick()
}

// Update handles terminal messages and advances the game.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.terminalWidth = msg.Width
		a.terminalHeight = msg.Height
		if a.engine != nil && (a.mode == "local" || a.mode == "host") {
			layout := ui.NewLayout(msg.Width, msg.Height)
			if layout.Playable {
				a.engine.Resize(layout.BoardWidth, layout.BoardHeight)
			}
		}
	case tea.KeyPressMsg:
		if a.mode == "join" {
			action := keyActionFor(msg)
			switch action {
			case keyActionQuit:
				return a, tea.Quit
			default:
				if direction, ok := directionForAction(action); ok && a.clientSession != nil {
					if err := a.clientSession.SendInput(direction); err != nil {
						a.err = err
					}
				}
			}

			return a, nil
		}

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
		if a.engine == nil {
			return a, tick()
		}

		if a.mode == "host" && a.hostSession != nil {
			if a.hostSession.Ready() {
				a.engine.SetPlayerCount(2)
			} else {
				a.engine.SetPlayerCount(1)
			}

			for {
				select {
				case input := <-a.hostSession.Inputs():
					a.engine.TurnPlayer(input.PlayerID, input.Direction)
				default:
					goto tickDone
				}
			}
		}

	tickDone:
		a.engine.Tick()
		if a.mode == "host" && a.hostSession != nil {
			if err := a.hostSession.SendSnapshot(a.engine.State()); err != nil {
				a.err = err
			}
		}
		return a, tick()
	case clientConnectedMsg:
		a.clientSession = msg.client
		a.err = msg.err
		if a.clientSession == nil {
			return a, nil
		}

		return a, pollState(func() (game.State, bool) {
			state, ok := <-a.clientSession.Snapshots()
			return state, ok
		})
	case stateMsg:
		a.remoteState = game.State(msg)
		return a, pollState(func() (game.State, bool) {
			if a.clientSession == nil {
				return game.State{}, false
			}
			state, ok := <-a.clientSession.Snapshots()
			return state, ok
		})
	}

	return a, nil
}

// View renders the current game screen.
func (a App) View() tea.View {
	if a.err != nil {
		view := tea.NewView(fmt.Sprintf("Mode: %s\nError: %v", a.mode, a.err))
		view.AltScreen = true
		view.WindowTitle = "Snake"

		return view
	}

	if a.mode == "join" && a.clientSession == nil {
		view := tea.NewView(fmt.Sprintf("Connecting to host %s...", a.joinAddr))
		view.AltScreen = true
		view.WindowTitle = "Snake"

		return view
	}

	if a.mode == "join" && a.remoteState.Width == 0 {
		view := tea.NewView(
			fmt.Sprintf(
				"Connected to %s as player %d.\nWaiting for first snapshot...",
				a.joinAddr,
				a.clientSession.PlayerID(),
			),
		)
		view.AltScreen = true
		view.WindowTitle = "Snake"

		return view
	}

	state := game.State{}
	switch {
	case a.mode == "join" && a.remoteState.Width > 0:
		state = a.remoteState
	case a.engine != nil:
		state = a.engine.State()
	}

	content := a.renderer.RenderRoom(state, a.terminalWidth, a.terminalHeight, a.roomMembers(state))

	view := tea.NewView(content)
	view.AltScreen = true
	view.WindowTitle = "Snake"

	return view
}

func (a App) roomMembers(state game.State) []ui.RoomMember {
	switch a.mode {
	case "host":
		if a.hostSession == nil {
			return nil
		}

		members := []ui.RoomMember{
			{
				Name:    "Player 1",
				Addr:    a.hostSession.Addr(),
				Status:  "hosting",
				Current: true,
			},
		}

		if a.hostSession.Ready() {
			members = append(members, ui.RoomMember{
				Name:      "Player 2",
				Addr:      a.hostSession.RemoteAddr(),
				Status:    "joined",
				LatencyMS: a.hostSession.LatencyMS(),
			})
		} else {
			members = append(members, ui.RoomMember{
				Name:   "Player 2",
				Status: "waiting",
			})
		}

		return members
	case "join":
		if a.clientSession == nil {
			return nil
		}

		members := []ui.RoomMember{
			{
				Name:   "Host",
				Addr:   a.joinAddr,
				Status: "connected",
			},
		}

		status := "joined"
		if state.Width == 0 {
			status = "waiting for snapshot"
		}

		members = append(members, ui.RoomMember{
			Name:      fmt.Sprintf("Player %d", a.clientSession.PlayerID()),
			Addr:      a.clientSession.LocalAddr(),
			Status:    status,
			LatencyMS: a.clientSession.LatencyMS(),
			Current:   true,
		})

		return members
	default:
		return nil
	}
}
