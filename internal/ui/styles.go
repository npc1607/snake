package ui

import "charm.land/lipgloss/v2"

type styles struct {
	frame       lipgloss.Style
	title       lipgloss.Style
	status      lipgloss.Style
	help        lipgloss.Style
	board       lipgloss.Style
	sidebar     lipgloss.Style
	sidebarHead lipgloss.Style
	rank        lipgloss.Style
	rankTop1    lipgloss.Style
	rankTop2    lipgloss.Style
	rankTop3    lipgloss.Style
	emptyRank   lipgloss.Style
	emptyEven   lipgloss.Style
	emptyOdd    lipgloss.Style
	food        lipgloss.Style
	snakeHead   lipgloss.Style
	message     lipgloss.Style
	smallScreen lipgloss.Style
}

func newStyles() styles {
	return styles{
		frame: lipgloss.NewStyle().
			Foreground(textColor).
			Background(appBackground).
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor),
		status: lipgloss.NewStyle().
			Foreground(textColor),
		help: lipgloss.NewStyle().
			Foreground(mutedTextColor),
		board: lipgloss.NewStyle().
			Background(boardBackground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor),
		sidebar: lipgloss.NewStyle().
			Background(boardBackground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 1),
		sidebarHead: lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor),
		rank: lipgloss.NewStyle().
			Foreground(textColor),
		rankTop1: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FBBF24")),
		rankTop2: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#CBD5E1")),
		rankTop3: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FB923C")),
		emptyRank: lipgloss.NewStyle().
			Foreground(mutedTextColor),
		emptyEven: lipgloss.NewStyle().
			Background(cellEven),
		emptyOdd: lipgloss.NewStyle().
			Background(cellOdd),
		food: lipgloss.NewStyle().
			Background(foodColor).
			Foreground(foodColor),
		snakeHead: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#E0F2FE")).
			Foreground(lipgloss.Color("#E0F2FE")),
		message: lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true),
		smallScreen: lipgloss.NewStyle().
			Foreground(textColor).
			Background(appBackground).
			Padding(1, 2),
	}
}
