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
	engine.players[0].snake = []Point{{X: 11, Y: 3}, {X: 10, Y: 3}, {X: 9, Y: 3}}
	engine.players[0].direction = DirectionRight
	engine.players[0].next = DirectionRight
	engine.players[0].alive = true
	engine.players[0].started = true
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

func TestSecondPlayerDoesNotMoveBeforeInput(t *testing.T) {
	engine := NewEngine(Config{Width: 20, Height: 12, PlayerCount: 2}, rand.New(rand.NewSource(1)))
	engine.Turn(DirectionRight)

	before := engine.State()
	if len(before.Players) != 2 {
		t.Fatalf("player count = %d, want 2", len(before.Players))
	}
	if before.Players[1].Alive {
		t.Fatalf("player 2 should be inactive before input")
	}

	engine.Tick()
	after := engine.State()

	if after.Players[1].Alive {
		t.Fatalf("player 2 should still be inactive without input")
	}
	if after.Players[1].Snake[0] != before.Players[1].Snake[0] {
		t.Fatalf("player 2 head moved: got %+v, want %+v", after.Players[1].Snake[0], before.Players[1].Snake[0])
	}
}

func TestFirstPlayerDoesNotMoveWhenSecondPlayerStarts(t *testing.T) {
	engine := NewEngine(Config{Width: 20, Height: 12, PlayerCount: 2}, rand.New(rand.NewSource(1)))

	before := engine.State()
	engine.TurnPlayer(2, DirectionUp)
	engine.Tick()
	after := engine.State()

	if after.Status != StatusRunning {
		t.Fatalf("status = %v, want %v", after.Status, StatusRunning)
	}
	if !after.Players[1].Alive {
		t.Fatalf("player 2 should become active after input")
	}
	if after.Players[0].Snake[0] != before.Players[0].Snake[0] {
		t.Fatalf("player 1 head moved: got %+v, want %+v", after.Players[0].Snake[0], before.Players[0].Snake[0])
	}

	wantPlayerTwoHead := Point{X: before.Players[1].Snake[0].X, Y: before.Players[1].Snake[0].Y - 1}
	if after.Players[1].Snake[0] != wantPlayerTwoHead {
		t.Fatalf("player 2 head = %+v, want %+v", after.Players[1].Snake[0], wantPlayerTwoHead)
	}
}

func TestPausedSecondPlayerInputDoesNotActivateSnake(t *testing.T) {
	engine := NewEngine(Config{Width: 20, Height: 12, PlayerCount: 2}, rand.New(rand.NewSource(1)))
	engine.Turn(DirectionRight)
	engine.TogglePause()

	before := engine.State()
	engine.TurnPlayer(2, DirectionUp)
	after := engine.State()

	if after.Status != StatusPaused {
		t.Fatalf("status = %v, want %v", after.Status, StatusPaused)
	}
	if after.Players[1].Alive {
		t.Fatalf("player 2 should remain inactive while paused")
	}
	if after.Players[1].Snake[0] != before.Players[1].Snake[0] {
		t.Fatalf("player 2 head changed while paused: got %+v, want %+v", after.Players[1].Snake[0], before.Players[1].Snake[0])
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
