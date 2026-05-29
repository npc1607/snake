package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"snake/internal/tui"
)

func main() {
	mode, hostAddr, joinAddr, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		os.Exit(2)
	}

	app := tui.NewLocalApp()
	switch mode {
	case "host":
		app = tui.NewHostApp(hostAddr)
	case "join":
		app = tui.NewJoinApp(joinAddr)
	}

	program := tea.NewProgram(app)
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run snake: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags(args []string) (mode string, hostAddr string, joinAddr string, err error) {
	flags := flag.NewFlagSet("snake", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	hostMode := flags.Bool("host", false, "start a LAN host")
	hostAddrValue := flags.String("host-addr", ":7777", "host listen address")
	joinMode := flags.String("join", "", "join a LAN host address")
	if err := flags.Parse(args); err != nil {
		return "", "", "", err
	}

	if *hostMode && *joinMode != "" {
		return "", "", "", fmt.Errorf("cannot use --host and --join together")
	}

	switch {
	case *hostMode:
		return "host", *hostAddrValue, "", nil
	case *joinMode != "":
		return "join", "", *joinMode, nil
	default:
		return "local", "", "", nil
	}
}
