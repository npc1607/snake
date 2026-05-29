package game

import (
	"math/rand"

	"snake/internal/helper"
)

type playerSnake struct {
	id        int
	snake     []Point
	direction Direction
	next      Direction
	alive     bool
	started   bool
}

// Engine owns the game rules and mutable state.
type Engine struct {
	width                   int
	height                  int
	playerCount             int
	maxSnakes               int
	enemySpawnIntervalTicks int
	ticksSinceEnemy         int
	nextEnemyID             int
	players                 []playerSnake
	enemies                 []enemySnake
	drops                   []Point
	food                    Point
	score                   int
	recorded                bool
	scores                  []int
	frame                   int
	status                  Status
	random                  *rand.Rand
}

// NewEngine creates a game engine with the provided configuration.
func NewEngine(config Config, random *rand.Rand) *Engine {
	config.Width, config.Height = normalizeSize(config.Width, config.Height)
	if config.PlayerCount <= 0 {
		config.PlayerCount = 1
	}
	if config.PlayerCount > 2 {
		config.PlayerCount = 2
	}
	if config.MaxSnakes <= 1 {
		config.MaxSnakes = defaultMaxSnakes
	}

	if config.EnemySpawnIntervalTicks <= 0 {
		config.EnemySpawnIntervalTicks = defaultEnemySpawnIntervalTicks
	}

	if random == nil {
		random = rand.New(rand.NewSource(1))
	}

	engine := &Engine{
		width:                   config.Width,
		height:                  config.Height,
		playerCount:             config.PlayerCount,
		maxSnakes:               config.MaxSnakes,
		enemySpawnIntervalTicks: config.EnemySpawnIntervalTicks,
		random:                  random,
	}
	engine.Reset()

	return engine
}

// State returns a snapshot of the current game state.
func (e *Engine) State() State {
	players := playerStates(e.players)
	snake := []Point(nil)
	if len(players) > 0 {
		snake = make([]Point, len(players[0].Snake))
		copy(snake, players[0].Snake)
	}

	direction := DirectionRight
	if len(e.players) > 0 {
		direction = e.players[0].direction
	}

	return State{
		Width:     e.width,
		Height:    e.height,
		Snake:     snake,
		Players:   players,
		Enemies:   enemyStates(e.enemies),
		Drops:     helper.CopySlice(e.drops),
		Food:      e.food,
		Direction: direction,
		Score:     e.score,
		Scores:    scoreEntries(e.scores),
		Frame:     e.frame,
		Status:    e.status,
		MaxSnakes: e.maxSnakes,
		NextEnemy: e.nextEnemyCountdown(),
	}
}

// SetPlayerCount updates the active player count and starts a fresh game.
func (e *Engine) SetPlayerCount(playerCount int) {
	if playerCount <= 0 {
		playerCount = 1
	}
	if playerCount > 2 {
		playerCount = 2
	}
	if playerCount == e.playerCount {
		return
	}

	e.playerCount = playerCount
	e.Reset()
}

// Resize changes the board size and starts a fresh game.
func (e *Engine) Resize(width int, height int) {
	width, height = normalizeSize(width, height)
	if width == e.width && height == e.height {
		return
	}

	e.width = width
	e.height = height
	e.Reset()
}

// Reset starts a new game.
func (e *Engine) Reset() {
	e.recordScore()

	center := Point{X: e.width / 2, Y: e.height / 2}

	e.score = 0
	e.recorded = false
	e.frame = 0
	e.ticksSinceEnemy = 0
	e.status = StatusReady
	e.players = e.initialPlayers(center)
	e.enemies = nil
	e.drops = nil
	e.food = e.nextFood()
}

// Turn queues a direction change for the next tick.
func (e *Engine) Turn(direction Direction) {
	e.TurnPlayer(1, direction)
}

// TurnPlayer queues a direction change for one player on the next tick.
func (e *Engine) TurnPlayer(playerID int, direction Direction) {
	if e.status == StatusGameOver {
		return
	}

	playerIndex := playerID - 1
	if playerIndex < 0 || playerIndex >= len(e.players) {
		return
	}

	player := &e.players[playerIndex]
	if player.direction.IsOpposite(direction) {
		return
	}

	if !player.alive {
		if e.status == StatusPaused {
			return
		}
		player.alive = true
	}

	player.started = true
	player.next = direction
	if e.status == StatusReady {
		e.status = StatusRunning
	}
}

// TogglePause pauses a running game or resumes a paused game.
func (e *Engine) TogglePause() {
	switch e.status {
	case StatusRunning:
		e.status = StatusPaused
	case StatusPaused, StatusReady:
		e.status = StatusRunning
	}
}

// Tick advances the game by one frame.
func (e *Engine) Tick() {
	e.frame++

	if e.status != StatusRunning {
		return
	}

	e.maybeSpawnEnemy()
	e.stepPlayers()
	if e.status != StatusRunning {
		return
	}
	e.moveEnemies()
}

func normalizeSize(width int, height int) (int, int) {
	if width < 12 {
		width = 12
	}

	if height < 8 {
		height = 8
	}

	return width, height
}

func (e *Engine) hitsWall(point Point) bool {
	return point.X < 0 || point.X >= e.width || point.Y < 0 || point.Y >= e.height
}

func (e *Engine) hitsPlayerSelf(playerIndex int, point Point, grows bool) bool {
	body := e.players[playerIndex].snake
	if !grows {
		body = body[:len(body)-1]
	}

	for _, segment := range body {
		if segment == point {
			return true
		}
	}

	return false
}

func (e *Engine) initialPlayers(center Point) []playerSnake {
	players := make([]playerSnake, 0, e.playerCount)
	for index := 0; index < e.playerCount; index++ {
		head := center
		direction := DirectionRight
		switch index {
		case 0:
			head = center
			direction = DirectionRight
		case 1:
			head = Point{X: center.X, Y: center.Y + 3}
			if head.Y >= e.height {
				head.Y = center.Y - 3
			}
			direction = DirectionRight
		}

		snake := []Point{
			head,
			{X: head.X - 1, Y: head.Y},
			{X: head.X - 2, Y: head.Y},
			{X: head.X - 3, Y: head.Y},
		}

		players = append(players, playerSnake{
			id:        index + 1,
			snake:     snake,
			direction: direction,
			next:      direction,
			alive:     index == 0,
			started:   false,
		})
	}

	return players
}

func (e *Engine) stepPlayers() {
	type playerMove struct {
		head      Point
		grows     bool
		ateFood   bool
		dropIndex int
	}

	moves := make([]playerMove, len(e.players))
	headCounts := make(map[Point]int, len(e.players))

	for index := range e.players {
		player := &e.players[index]
		if !player.alive || !player.started || len(player.snake) == 0 {
			continue
		}

		player.direction = player.next
		head := player.snake[0].move(player.direction)
		ateFood := head == e.food
		dropIndex := e.dropIndex(head)

		moves[index] = playerMove{
			head:      head,
			grows:     ateFood || dropIndex >= 0,
			ateFood:   ateFood,
			dropIndex: dropIndex,
		}
		headCounts[head]++
	}

	for index := range e.players {
		player := &e.players[index]
		if !player.alive || !player.started || len(player.snake) == 0 {
			continue
		}

		move := moves[index]
		if e.hitsWall(move.head) || e.hitsEnemy(move.head) || e.hitsPlayerSelf(index, move.head, move.grows) {
			e.status = StatusGameOver
			e.recordScore()
			return
		}

		if headCounts[move.head] > 1 {
			e.status = StatusGameOver
			e.recordScore()
			return
		}

		for otherIndex := range e.players {
			if otherIndex == index {
				continue
			}

			other := e.players[otherIndex]
			if !other.alive || len(other.snake) == 0 {
				continue
			}

			if move.head == other.snake[0] && moves[otherIndex].head == player.snake[0] {
				e.status = StatusGameOver
				e.recordScore()
				return
			}

			if e.hitsOtherPlayerBody(otherIndex, move.head, moves[otherIndex].grows) {
				e.status = StatusGameOver
				e.recordScore()
				return
			}
		}
	}

	foodEaten := false
	dropsToRemove := make(map[int]struct{}, len(e.players))
	for index := range e.players {
		player := &e.players[index]
		if !player.alive || !player.started || len(player.snake) == 0 {
			continue
		}

		move := moves[index]
		nextSnake := make([]Point, 0, len(player.snake)+1)
		nextSnake = append(nextSnake, move.head)
		nextSnake = append(nextSnake, player.snake...)
		if !move.grows {
			nextSnake = nextSnake[:len(nextSnake)-1]
		}
		player.snake = nextSnake

		if move.ateFood || move.dropIndex >= 0 {
			e.score++
			e.recorded = false
		}
		if move.ateFood {
			foodEaten = true
		}
		if move.dropIndex >= 0 {
			dropsToRemove[move.dropIndex] = struct{}{}
		}
	}

	if len(dropsToRemove) > 0 {
		e.removeDrops(dropsToRemove)
	}
	if foodEaten {
		e.food = e.nextFood()
	}
}

func (e *Engine) hitsOtherPlayerBody(playerIndex int, point Point, otherGrows bool) bool {
	body := e.players[playerIndex].snake
	if len(body) == 0 {
		return false
	}
	limit := len(body)
	if !otherGrows {
		limit--
	}
	if limit < 1 {
		limit = 1
	}
	for segmentIndex := 1; segmentIndex < limit; segmentIndex++ {
		if body[segmentIndex] == point {
			return true
		}
	}

	return false
}

func (e *Engine) removeDrops(indexes map[int]struct{}) {
	if len(indexes) == 0 {
		return
	}

	nextDrops := e.drops[:0]
	for index, drop := range e.drops {
		if _, ok := indexes[index]; ok {
			continue
		}
		nextDrops = append(nextDrops, drop)
	}

	e.drops = nextDrops
}

func (e *Engine) nextFood() Point {
	occupied := e.occupiedCells()

	freeCells := e.width*e.height - len(occupied)
	if freeCells <= 0 {
		return Point{}
	}

	target := e.random.Intn(freeCells)
	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {
			point := Point{X: x, Y: y}
			if _, exists := occupied[point]; exists {
				continue
			}

			if target == 0 {
				return point
			}
			target--
		}
	}

	return Point{}
}
