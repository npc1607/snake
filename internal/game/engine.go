package game

import "math/rand"

// Engine owns the game rules and mutable state.
type Engine struct {
	width     int
	height    int
	snake     []Point
	food      Point
	direction Direction
	next      Direction
	score     int
	recorded  bool
	scores    []int
	frame     int
	status    Status
	random    *rand.Rand
}

// NewEngine creates a game engine with the provided configuration.
func NewEngine(config Config, random *rand.Rand) *Engine {
	config.Width, config.Height = normalizeSize(config.Width, config.Height)

	if random == nil {
		random = rand.New(rand.NewSource(1))
	}

	engine := &Engine{
		width:  config.Width,
		height: config.Height,
		random: random,
	}
	engine.Reset()

	return engine
}

// State returns a snapshot of the current game state.
func (e *Engine) State() State {
	snake := make([]Point, len(e.snake))
	copy(snake, e.snake)

	return State{
		Width:     e.width,
		Height:    e.height,
		Snake:     snake,
		Food:      e.food,
		Direction: e.direction,
		Score:     e.score,
		Scores:    scoreEntries(e.scores),
		Frame:     e.frame,
		Status:    e.status,
	}
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

	e.direction = DirectionRight
	e.next = DirectionRight
	e.score = 0
	e.recorded = false
	e.frame = 0
	e.status = StatusReady
	e.snake = []Point{
		center,
		{X: center.X - 1, Y: center.Y},
		{X: center.X - 2, Y: center.Y},
		{X: center.X - 3, Y: center.Y},
	}
	e.food = e.nextFood()
}

// Turn queues a direction change for the next tick.
func (e *Engine) Turn(direction Direction) {
	if e.status == StatusGameOver {
		return
	}

	if e.direction.IsOpposite(direction) {
		return
	}

	e.next = direction
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

	e.direction = e.next
	head := e.snake[0].move(e.direction)
	ateFood := head == e.food

	if e.hitsWall(head) || e.hitsSelf(head, ateFood) {
		e.status = StatusGameOver
		e.recordScore()
		return
	}

	nextSnake := make([]Point, 0, len(e.snake)+1)
	nextSnake = append(nextSnake, head)
	nextSnake = append(nextSnake, e.snake...)

	if ateFood {
		e.score++
		e.recorded = false
		e.snake = nextSnake
		e.food = e.nextFood()
		return
	}

	e.snake = nextSnake[:len(nextSnake)-1]
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

func (e *Engine) hitsSelf(point Point, grows bool) bool {
	body := e.snake
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

func (e *Engine) nextFood() Point {
	occupied := make(map[Point]struct{}, len(e.snake))
	for _, segment := range e.snake {
		occupied[segment] = struct{}{}
	}

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
