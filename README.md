# Snake

Snake is a terminal snake game written in Go with [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) and [Lip Gloss v2](https://github.com/charmbracelet/lipgloss). It keeps the classic snake loop, then adds adaptive terminal layout, animated colors, enemy snake AI, loot drops, LAN multiplayer, and an in-game Top 10 score board.

## Features

- Adaptive terminal board that fills the available window while keeping a right-side status panel.
- Color-shifting player snake body.
- Top 10 leaderboard ranked by eating score.
- Enemy snakes that spawn over time and use BFS pathfinding to block food or attack the player.
- Enemy death drops green loot blocks that the player can eat for points and growth.
- 2-player LAN host/join mode with a shared board.
- Right-side room member list showing player address, join state, and latency.
- Clean internal layering so game rules, Bubble Tea event handling, and Lip Gloss rendering stay separate.

## Requirements

- Go 1.26 or newer.
- A terminal with ANSI color support.

## Run

```bash
make run
```

Or run directly:

```bash
go run ./cmd/snake
```

LAN multiplayer modes:

```bash
go run ./cmd/snake --host
go run ./cmd/snake --host --host-addr 0.0.0.0:7777
go run ./cmd/snake --join 192.168.1.10:7777
```

`--host` listens on `:7777` by default. Use `--host-addr` to bind another interface or port.

In LAN mode, the right sidebar shows a `ROOM` section with the current players:

- Host view: `Player 1` is the host, `Player 2` is `waiting` or `joined`.
- Join view: the host address and the local player are shown.
- Latency is measured with lightweight ping/pong messages. Latency below 60ms is shown in green; latency of 60ms or higher is shown in red.

## Controls

| Action | Keys |
| --- | --- |
| Move | Arrow keys, `WASD`, `HJKL` |
| Pause / resume | `Space`, `P` |
| Restart | `R` |
| Quit | `Q`, `Esc`, `Ctrl+C` |

## Gameplay

Eat food to gain score and grow. Every 30 seconds, the game attempts to spawn a new enemy snake. The board can hold up to 6 snakes total: the player plus 5 enemies.

Enemy snakes use pathfinding to either chase food or pressure the player. If an enemy hits the player's head, the game is over. If an enemy hits the player's body, a wall, or another enemy, it breaks apart and leaves green loot blocks. Eating those blocks also increases score and length.

## Development

```bash
make help    # Show available targets.
make test    # Run all tests.
make build   # Build bin/snake.
make fmt     # Format Go files.
make tidy    # Tidy module dependencies.
make clean   # Remove build artifacts.
```

## Project Structure

```text
cmd/snake
  Program entry point.

internal/game
  Pure game rules: movement, scoring, collisions, enemies, drops, leaderboard, pathfinding, and shared 2-player state.

internal/session
  LAN host/client sessions: join handshake, input forwarding, snapshots, and latency ping/pong.

internal/tui
  Bubble Tea adapter: local/host/join model lifecycle, key handling, ticks, and terminal resize events.

internal/ui
  Lip Gloss renderer: layout, colors, board drawing, sidebar, room members, and status text.

internal/helper
  Small reusable helpers for address network selection, slice copying, and text truncation.

docs/design.md
  Design notes and extension points.
```

## Architecture Notes

The core game engine does not import Bubble Tea or Lip Gloss. That keeps the domain logic testable and makes it possible to add other renderers or controllers later. The current AI uses breadth-first search over the board grid to choose enemy movement toward either the food or a player's head.

## LAN Multiplayer Roadmap

The LAN design is documented in [docs/design.md](/Users/pcong/Project/snake/docs/design.md:94).

Implemented:

1. 2-player shared-map LAN mode while keeping single-player as the default with no CLI flags.
2. Host/join CLI flow.
3. Room member display with IP/address, connection status, and latency.

Still planned:

1. Room discovery and room listing.
2. Reconnect support after disconnects.
