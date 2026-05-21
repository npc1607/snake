package ui

const (
	frameHorizontalPadding = 4
	frameVerticalPadding   = 2
	sidebarWidth           = 26
	layoutGap              = 2
	minBoardWidth          = 12
	minBoardHeight         = 8
	defaultTerminalWidth   = 100
	defaultTerminalHeight  = 30
)

// Layout describes the terminal space assigned to the game surface.
type Layout struct {
	TerminalWidth  int
	TerminalHeight int
	ContentWidth   int
	ContentHeight  int
	BoardWidth     int
	BoardHeight    int
	LeftWidth      int
	SidebarWidth   int
	Gap            int
	Playable       bool
}

// NewLayout calculates board and sidebar dimensions for the current terminal.
func NewLayout(terminalWidth int, terminalHeight int) Layout {
	if terminalWidth <= 0 {
		terminalWidth = defaultTerminalWidth
	}

	if terminalHeight <= 0 {
		terminalHeight = defaultTerminalHeight
	}

	contentWidth := terminalWidth - frameHorizontalPadding
	contentHeight := terminalHeight - frameVerticalPadding
	layout := Layout{
		TerminalWidth:  terminalWidth,
		TerminalHeight: terminalHeight,
		ContentWidth:   contentWidth,
		ContentHeight:  contentHeight,
		SidebarWidth:   sidebarWidth,
		Gap:            layoutGap,
	}

	boardWidth := (contentWidth - sidebarWidth - layoutGap - 2) / len(cellContent)
	boardHeight := contentHeight - 4
	if boardWidth < minBoardWidth || boardHeight < minBoardHeight {
		return layout
	}

	layout.BoardWidth = boardWidth
	layout.BoardHeight = boardHeight
	layout.LeftWidth = boardWidth*len(cellContent) + 2
	layout.Playable = true

	return layout
}
