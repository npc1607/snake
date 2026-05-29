package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var snakeWave = []color.Color{
	lipgloss.Color("#34D399"),
	lipgloss.Color("#22D3EE"),
	lipgloss.Color("#60A5FA"),
	lipgloss.Color("#A78BFA"),
	lipgloss.Color("#F472B6"),
	lipgloss.Color("#FBBF24"),
}

var snakeWaveAlt = []color.Color{
	lipgloss.Color("#F59E0B"),
	lipgloss.Color("#F97316"),
	lipgloss.Color("#FB7185"),
	lipgloss.Color("#F43F5E"),
	lipgloss.Color("#E879F9"),
	lipgloss.Color("#C084FC"),
}

var enemyWave = []color.Color{
	lipgloss.Color("#EF4444"),
	lipgloss.Color("#F97316"),
	lipgloss.Color("#F59E0B"),
	lipgloss.Color("#DC2626"),
}

var (
	appBackground   = lipgloss.Color("#101418")
	boardBackground = lipgloss.Color("#121A20")
	cellEven        = lipgloss.Color("#17212A")
	cellOdd         = lipgloss.Color("#1B2630")
	foodColor       = lipgloss.Color("#F43F5E")
	dropColor       = lipgloss.Color("#A3E635")
	textColor       = lipgloss.Color("#D8DEE9")
	mutedTextColor  = lipgloss.Color("#7F8C98")
	accentColor     = lipgloss.Color("#38BDF8")
	warningColor    = lipgloss.Color("#F59E0B")
	dangerColor     = lipgloss.Color("#FB7185")
	borderColor     = lipgloss.Color("#2D3A45")
)
