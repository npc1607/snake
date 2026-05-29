package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"snake/internal/game"
	"snake/internal/helper"
)

const cellContent = "  "

// Renderer draws the terminal UI for the game.
type Renderer struct {
	styles styles
}

// RoomMember describes one visible LAN room participant.
type RoomMember struct {
	Name      string
	Addr      string
	Status    string
	LatencyMS int
	Current   bool
}

// NewRenderer creates the default game renderer.
func NewRenderer() Renderer {
	return Renderer{
		styles: newStyles(),
	}
}

// Render converts the game state into styled terminal content.
func (r Renderer) Render(state game.State, terminalWidth int, terminalHeight int) string {
	return r.RenderRoom(state, terminalWidth, terminalHeight, nil)
}

// RenderRoom converts the game state and room members into styled terminal content.
func (r Renderer) RenderRoom(
	state game.State,
	terminalWidth int,
	terminalHeight int,
	roomMembers []RoomMember,
) string {
	layout := NewLayout(terminalWidth, terminalHeight)
	if !layout.Playable {
		return r.renderSmallScreen(layout)
	}

	header := r.renderHeader(state, layout.LeftWidth)
	board := r.renderBoard(state)
	help := r.renderHelp(layout.LeftWidth)
	left := lipgloss.JoinVertical(lipgloss.Left, header, board, help)
	sidebar := r.renderSidebar(state, layout, roomMembers)
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", layout.Gap), sidebar)
	content = lipgloss.Place(layout.ContentWidth, layout.ContentHeight, lipgloss.Left, lipgloss.Top, content)

	return r.styles.frame.Render(content)
}

func (r Renderer) renderHeader(state game.State, width int) string {
	title := r.styles.title.Render("SNAKE")
	playerCount := len(state.Players)
	if playerCount == 0 {
		playerCount = 1
	}
	maxEnemies := state.MaxSnakes - playerCount
	if maxEnemies < 0 {
		maxEnemies = 0
	}
	status := r.styles.status.Render(fmt.Sprintf(
		"Score: %d  Players: %d  Enemies: %d/%d  State: %s",
		state.Score,
		playerCount,
		len(state.Enemies),
		maxEnemies,
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
	primarySnake := state.Snake
	if len(state.Players) > 0 {
		primarySnake = state.Players[0].Snake
	}
	for index, point := range primarySnake {
		snakeCells[point] = index
	}

	secondSnakeCells := make(map[game.Point]int)
	if len(state.Players) > 1 && state.Players[1].Alive {
		for index, point := range state.Players[1].Snake {
			secondSnakeCells[point] = index
		}
	}

	enemyCells := make(map[game.Point]int)
	enemyHeads := make(map[game.Point]struct{}, len(state.Enemies))
	for _, enemy := range state.Enemies {
		for index, point := range enemy.Snake {
			if index == 0 {
				enemyHeads[point] = struct{}{}
			}
			enemyCells[point] = index
		}
	}

	dropCells := make(map[game.Point]struct{}, len(state.Drops))
	for _, drop := range state.Drops {
		dropCells[drop] = struct{}{}
	}

	rows := make([]string, 0, state.Height)
	for y := 0; y < state.Height; y++ {
		var row strings.Builder
		for x := 0; x < state.Width; x++ {
			point := game.Point{X: x, Y: y}
			row.WriteString(
				r.renderCell(state, point, snakeCells, secondSnakeCells, enemyCells, enemyHeads, dropCells),
			)
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
	help := "Move: arrows/WASD/HJKL  Pause: space/P  Restart: R  Quit: Q/Esc  Lure enemies into your body"
	if lipgloss.Width(help) > width {
		help = "Move: arrows/WASD  Pause: space  Restart: R  Quit: Q  Lure enemies"
	}

	return r.styles.help.Width(width).Render(help)
}

func (r Renderer) renderSidebar(state game.State, layout Layout, roomMembers []RoomMember) string {
	playerCount := len(state.Players)
	if playerCount == 0 {
		playerCount = 1
	}
	maxEnemies := state.MaxSnakes - playerCount
	if maxEnemies < 0 {
		maxEnemies = 0
	}

	rows := []string{
		r.styles.sidebarHead.Render("TOP 10"),
		r.styles.help.Render("Best eating scores"),
		"",
	}

	for rank := 1; rank <= 10; rank++ {
		rows = append(rows, r.renderRank(state.Scores, rank))
	}

	if len(roomMembers) > 0 {
		rows = append(rows, "", r.styles.sidebarHead.Render("ROOM"))
		for _, member := range roomMembers {
			rows = append(rows, r.renderRoomMember(member, layout.SidebarWidth-4))
		}
	}

	rows = append(rows,
		"",
		r.styles.sidebarHead.Render("THREATS"),
		r.styles.rank.Render(fmt.Sprintf("Enemies  %d/%d", len(state.Enemies), maxEnemies)),
		r.styles.rank.Render(fmt.Sprintf("Next     %ds", enemySeconds(state.NextEnemy))),
		r.styles.help.Render("Green blocks are loot."),
	)

	content := strings.Join(rows, "\n")

	return r.styles.sidebar.
		Width(layout.SidebarWidth - 2).
		Height(layout.ContentHeight - 2).
		Render(content)
}

func (r Renderer) renderRoomMember(member RoomMember, width int) string {
	name := member.Name
	if member.Current {
		name += " *"
	}

	if width < 1 {
		width = 1
	}

	address := member.Addr
	if address == "" {
		address = "-"
	}

	lines := []string{
		r.styles.rank.Render(name),
		r.styles.help.Render(helper.TruncateText(address, width)),
		r.styles.help.Render(helper.TruncateText(member.Status, width)),
	}
	if member.LatencyMS > 0 {
		lines = append(lines, r.renderLatency(member.LatencyMS, width))
	}

	return strings.Join(lines, "\n")
}

func (r Renderer) renderLatency(latencyMS int, width int) string {
	latencyText := helper.TruncateText(fmt.Sprintf("Latency %dms", latencyMS), width)
	style := r.styles.help
	if latencyMS > 0 && latencyMS < 60 {
		style = lipgloss.NewStyle().Foreground(dropColor)
	}
	if latencyMS >= 60 {
		style = lipgloss.NewStyle().Foreground(dangerColor)
	}

	return style.Render(latencyText)
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

func (r Renderer) renderCell(
	state game.State,
	point game.Point,
	snakeCells map[game.Point]int,
	secondSnakeCells map[game.Point]int,
	enemyCells map[game.Point]int,
	enemyHeads map[game.Point]struct{},
	dropCells map[game.Point]struct{},
) string {
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

	if index, ok := secondSnakeCells[point]; ok {
		if index == 0 {
			return r.styles.enemyHead.Render(cellContent)
		}

		color := snakeWaveAlt[(state.Frame+index)%len(snakeWaveAlt)]
		return lipgloss.NewStyle().
			Background(color).
			Foreground(color).
			Render(cellContent)
	}

	if _, ok := enemyHeads[point]; ok {
		return r.styles.enemyHead.Render(cellContent)
	}

	if index, ok := enemyCells[point]; ok {
		color := enemyWave[(state.Frame+index)%len(enemyWave)]
		return lipgloss.NewStyle().
			Background(color).
			Foreground(color).
			Render(cellContent)
	}

	if point == state.Food {
		return r.styles.food.Render(cellContent)
	}

	if _, ok := dropCells[point]; ok {
		return r.styles.drop.Render(cellContent)
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

func enemySeconds(ticks int) int {
	return (ticks*95 + 999) / 1000
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
