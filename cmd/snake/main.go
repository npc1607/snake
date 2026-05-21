package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"snake/internal/tui"
)

func main() {
	program := tea.NewProgram(tui.NewApp())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run snake: %v\n", err)
		os.Exit(1)
	}
}
