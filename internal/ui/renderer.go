package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"snake/internal/game"
)

const cellContent = "  "

// Renderer draws the terminal UI for the game.
type Renderer struct {
	styles styles
}

// NewRenderer creates the default game renderer.
func NewRenderer() Renderer {
	return Renderer{
		styles: newStyles(),
	}
}

// Render converts the game state into styled terminal content.
func (r Renderer) Render(state game.State, terminalWidth int, terminalHeight int) string {
	layout := NewLayout(terminalWidth, terminalHeight)
	if !layout.Playable {
		return r.renderSmallScreen(layout)
	}

	header := r.renderHeader(state, layout.LeftWidth)
	board := r.renderBoard(state)
	help := r.renderHelp(layout.LeftWidth)
	left := lipgloss.JoinVertical(lipgloss.Left, header, board, help)
	sidebar := r.renderLeaderboard(state.Scores, layout)
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", layout.Gap), sidebar)
	content = lipgloss.Place(layout.ContentWidth, layout.ContentHeight, lipgloss.Left, lipgloss.Top, content)

	return r.styles.frame.Render(content)
}

func (r Renderer) renderHeader(state game.State, width int) string {
	title := r.styles.title.Render("SNAKE")
	status := r.styles.status.Render(fmt.Sprintf(
		"Score: %d  Length: %d  State: %s",
		state.Score,
		len(state.Snake),
		statusText(state.Status),
	))
	spacerWidth := width - lipgloss.Width(title) - lipgloss.Width(status)
	if spacerWidth < 1 {
		spacerWidth = 1
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, strings.Repeat(" ", spacerWidth), status)
}

func (r Renderer) renderBoard(state game.State) string {
	snakeCells := make(map[game.Point]int, len(state.Snake))
	for index, point := range state.Snake {
		snakeCells[point] = index
	}

	rows := make([]string, 0, state.Height)
	for y := 0; y < state.Height; y++ {
		var row strings.Builder
		for x := 0; x < state.Width; x++ {
			point := game.Point{X: x, Y: y}
			row.WriteString(r.renderCell(state, point, snakeCells))
		}
		rows = append(rows, row.String())
	}

	board := strings.Join(rows, "\n")
	if message := overlayMessage(state.Status); message != "" {
		board = r.overlayCenteredMessage(board, state.Width, state.Height, message)
	}

	return r.styles.board.Render(board)
}

func (r Renderer) renderHelp(width int) string {
	help := "Move: arrows/WASD/HJKL  Pause: space/P  Restart: R  Quit: Q/Esc"
	if lipgloss.Width(help) > width {
		help = "Move: arrows/WASD  Pause: space  Restart: R  Quit: Q"
	}

	return r.styles.help.Width(width).Render(help)
}

func (r Renderer) renderLeaderboard(scores []game.ScoreEntry, layout Layout) string {
	rows := []string{
		r.styles.sidebarHead.Render("TOP 10"),
		r.styles.help.Render("Best food scores"),
		"",
	}

	for rank := 1; rank <= 10; rank++ {
		rows = append(rows, r.renderRank(scores, rank))
	}

	content := strings.Join(rows, "\n")

	return r.styles.sidebar.
		Width(layout.SidebarWidth - 2).
		Height(layout.ContentHeight - 2).
		Render(content)
}

func (r Renderer) renderRank(scores []game.ScoreEntry, rank int) string {
	marker := rankMarker(rank)
	scoreText := "--"
	if rank <= len(scores) {
		scoreText = fmt.Sprintf("%d", scores[rank-1].Score)
	}

	line := fmt.Sprintf("%-4s %2d  %8s", marker, rank, scoreText)
	switch rank {
	case 1:
		return r.styles.rankTop1.Render(line)
	case 2:
		return r.styles.rankTop2.Render(line)
	case 3:
		return r.styles.rankTop3.Render(line)
	default:
		if rank > len(scores) {
			return r.styles.emptyRank.Render(line)
		}

		return r.styles.rank.Render(line)
	}
}

func (r Renderer) renderCell(state game.State, point game.Point, snakeCells map[game.Point]int) string {
	if index, ok := snakeCells[point]; ok {
		if index == 0 {
			return r.styles.snakeHead.Render(cellContent)
		}

		color := snakeWave[(state.Frame+index)%len(snakeWave)]
		return lipgloss.NewStyle().
			Background(color).
			Foreground(color).
			Render(cellContent)
	}

	if point == state.Food {
		return r.styles.food.Render(cellContent)
	}

	if (point.X+point.Y)%2 == 0 {
		return r.styles.emptyEven.Render(cellContent)
	}

	return r.styles.emptyOdd.Render(cellContent)
}

func (r Renderer) overlayCenteredMessage(board string, width int, height int, message string) string {
	lines := strings.Split(board, "\n")
	if height == 0 || len(lines) == 0 {
		return board
	}

	message = " " + message + " "
	messageWidth := lipgloss.Width(message)
	boardWidth := width * len(cellContent)
	if messageWidth >= boardWidth {
		return board
	}

	row := height / 2
	leftPadding := (boardWidth - messageWidth) / 2
	rightPadding := boardWidth - messageWidth - leftPadding
	rendered := strings.Repeat(" ", leftPadding) +
		r.styles.message.Background(boardBackground).Render(message) +
		strings.Repeat(" ", rightPadding)

	lines[row] = rendered

	return strings.Join(lines, "\n")
}

func (r Renderer) renderSmallScreen(layout Layout) string {
	message := fmt.Sprintf(
		"Terminal too small. Need at least %dx%d, current size is %dx%d.",
		frameHorizontalPadding+sidebarWidth+layoutGap+minBoardWidth*len(cellContent)+2,
		frameVerticalPadding+minBoardHeight+4,
		layout.TerminalWidth,
		layout.TerminalHeight,
	)

	return r.styles.smallScreen.Render(message)
}

func rankMarker(rank int) string {
	switch rank {
	case 1:
		return "#1"
	case 2:
		return "#2"
	case 3:
		return "#3"
	default:
		return "--"
	}
}

func overlayMessage(status game.Status) string {
	switch status {
	case game.StatusReady:
		return "Press an arrow key to start"
	case game.StatusPaused:
		return "Paused"
	case game.StatusGameOver:
		return "Game over - press R"
	default:
		return ""
	}
}

func statusText(status game.Status) string {
	switch status {
	case game.StatusReady:
		return "ready"
	case game.StatusRunning:
		return "running"
	case game.StatusPaused:
		return "paused"
	case game.StatusGameOver:
		return "game over"
	default:
		return "unknown"
	}
}
