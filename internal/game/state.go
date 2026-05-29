package game

// Status describes the current game lifecycle state.
type Status int

const (
	StatusReady Status = iota
	StatusRunning
	StatusPaused
	StatusGameOver
)

// Config controls the game rules and board size.
type Config struct {
	Width                   int
	Height                  int
	PlayerCount             int
	MaxSnakes               int
	EnemySpawnIntervalTicks int
}

// DefaultConfig returns a balanced board size for terminal play.
func DefaultConfig() Config {
	return Config{
		Width:                   32,
		Height:                  20,
		PlayerCount:             1,
		MaxSnakes:               defaultMaxSnakes,
		EnemySpawnIntervalTicks: defaultEnemySpawnIntervalTicks,
	}
}

// PlayerState is a read-only player snake snapshot.
type PlayerState struct {
	ID        int
	Snake     []Point
	Direction Direction
	Alive     bool
}

// State is a read-only snapshot of the current game state for callers.
type State struct {
	Width     int
	Height    int
	Snake     []Point
	Players   []PlayerState
	Enemies   []EnemyState
	Drops     []Point
	Food      Point
	Direction Direction
	Score     int
	Scores    []ScoreEntry
	Frame     int
	Status    Status
	MaxSnakes int
	NextEnemy int
}
