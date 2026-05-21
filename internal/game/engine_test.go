package game

import (
	"math/rand"
	"testing"
)

func TestTickMovesSnake(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))

	before := engine.State()
	engine.Turn(DirectionRight)
	engine.Tick()
	after := engine.State()

	wantHead := Point{X: before.Snake[0].X + 1, Y: before.Snake[0].Y}
	if after.Snake[0] != wantHead {
		t.Fatalf("head = %+v, want %+v", after.Snake[0], wantHead)
	}

	if len(after.Snake) != len(before.Snake) {
		t.Fatalf("snake length = %d, want %d", len(after.Snake), len(before.Snake))
	}
}

func TestTickGrowsWhenEatingFood(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	state := engine.State()
	engine.food = Point{X: state.Snake[0].X + 1, Y: state.Snake[0].Y}

	engine.Turn(DirectionRight)
	engine.Tick()
	after := engine.State()

	if after.Score != 1 {
		t.Fatalf("score = %d, want 1", after.Score)
	}

	if len(after.Snake) != len(state.Snake)+1 {
		t.Fatalf("snake length = %d, want %d", len(after.Snake), len(state.Snake)+1)
	}
}

func TestTickGameOverWhenHittingWall(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	engine.snake = []Point{{X: 11, Y: 3}, {X: 10, Y: 3}, {X: 9, Y: 3}}
	engine.direction = DirectionRight
	engine.next = DirectionRight
	engine.status = StatusRunning

	engine.Tick()
	after := engine.State()

	if after.Status != StatusGameOver {
		t.Fatalf("status = %v, want %v", after.Status, StatusGameOver)
	}
}

func TestOppositeTurnIsIgnored(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	engine.Turn(DirectionLeft)
	engine.Tick()
	after := engine.State()

	if after.Direction != DirectionRight {
		t.Fatalf("direction = %v, want %v", after.Direction, DirectionRight)
	}
}

func TestPauseStopsMovement(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	engine.Turn(DirectionRight)
	engine.TogglePause()

	before := engine.State()
	engine.Tick()
	after := engine.State()

	if after.Status != StatusPaused {
		t.Fatalf("status = %v, want %v", after.Status, StatusPaused)
	}

	if after.Snake[0] != before.Snake[0] {
		t.Fatalf("head moved while paused: got %+v, want %+v", after.Snake[0], before.Snake[0])
	}
}

func TestLeaderboardKeepsTopTenScores(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	for score := 1; score <= 12; score++ {
		engine.score = score
		engine.recorded = false
		engine.recordScore()
	}

	entries := engine.State().Scores
	if len(entries) != 10 {
		t.Fatalf("leaderboard length = %d, want 10", len(entries))
	}

	if entries[0].Score != 12 {
		t.Fatalf("top score = %d, want 12", entries[0].Score)
	}

	if entries[9].Score != 3 {
		t.Fatalf("last score = %d, want 3", entries[9].Score)
	}
}

func TestLeaderboardDoesNotRecordSameScoreTwice(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	engine.score = 5
	engine.recordScore()
	engine.recordScore()

	entries := engine.State().Scores
	if len(entries) != 1 {
		t.Fatalf("leaderboard length = %d, want 1", len(entries))
	}
}

func TestResizeUpdatesBoardAndKeepsScores(t *testing.T) {
	engine := NewEngine(Config{Width: 12, Height: 8}, rand.New(rand.NewSource(1)))
	engine.score = 4
	engine.Resize(30, 18)

	state := engine.State()
	if state.Width != 30 || state.Height != 18 {
		t.Fatalf("board size = %dx%d, want 30x18", state.Width, state.Height)
	}

	if len(state.Scores) != 1 || state.Scores[0].Score != 4 {
		t.Fatalf("scores = %+v, want one score of 4", state.Scores)
	}

	if state.Score != 0 {
		t.Fatalf("current score = %d, want 0", state.Score)
	}
}
