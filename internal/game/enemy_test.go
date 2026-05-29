package game

import (
	"math/rand"
	"testing"
)

func TestEnemySpawnsAfterInterval(t *testing.T) {
	engine := NewEngine(
		Config{Width: 24, Height: 16, EnemySpawnIntervalTicks: 2, MaxSnakes: 6},
		rand.New(rand.NewSource(1)),
	)
	engine.Turn(DirectionRight)

	engine.Tick()
	if len(engine.State().Enemies) != 0 {
		t.Fatalf("enemy count = %d, want 0", len(engine.State().Enemies))
	}

	engine.Tick()
	if len(engine.State().Enemies) != 1 {
		t.Fatalf("enemy count = %d, want 1", len(engine.State().Enemies))
	}
}

func TestEnemyCountRespectsMaxSnakes(t *testing.T) {
	engine := NewEngine(
		Config{Width: 32, Height: 20, EnemySpawnIntervalTicks: 1, MaxSnakes: 6},
		rand.New(rand.NewSource(1)),
	)
	engine.Turn(DirectionRight)

	for tickIndex := 0; tickIndex < 20; tickIndex++ {
		engine.Tick()
	}

	if len(engine.State().Enemies) > 5 {
		t.Fatalf("enemy count = %d, want <= 5", len(engine.State().Enemies))
	}
}

func TestEnemyCrashScattersDrops(t *testing.T) {
	engine := NewEngine(Config{Width: 16, Height: 12}, rand.New(rand.NewSource(1)))
	enemy := enemySnake{
		id:        1,
		direction: DirectionRight,
		snake: []Point{
			{X: 4, Y: 5},
			{X: 3, Y: 5},
			{X: 2, Y: 5},
		},
	}
	engine.enemies = []enemySnake{enemy}

	engine.scatterEnemy(enemy)
	engine.enemies = nil
	state := engine.State()

	if len(state.Enemies) != 0 {
		t.Fatalf("enemy count = %d, want 0", len(state.Enemies))
	}

	if len(state.Drops) != 3 {
		t.Fatalf("drop count = %d, want 3", len(state.Drops))
	}
}

func TestEnemyHittingPlayerHeadEndsGame(t *testing.T) {
	engine := NewEngine(Config{Width: 16, Height: 12}, rand.New(rand.NewSource(1)))
	engine.players[0].snake = []Point{
		{X: 6, Y: 6},
		{X: 5, Y: 6},
		{X: 4, Y: 6},
	}
	engine.enemies = []enemySnake{
		{
			id:        1,
			direction: DirectionRight,
			snake: []Point{
				{X: 6, Y: 5},
				{X: 5, Y: 5},
				{X: 4, Y: 5},
			},
		},
	}

	engine.moveEnemies()
	if engine.State().Status != StatusGameOver {
		t.Fatalf("status = %v, want %v", engine.State().Status, StatusGameOver)
	}
}

func TestEnemyEatingFoodGrows(t *testing.T) {
	engine := NewEngine(Config{Width: 16, Height: 12}, rand.New(rand.NewSource(1)))
	engine.players[0].snake = []Point{
		{X: 10, Y: 10},
		{X: 9, Y: 10},
		{X: 8, Y: 10},
	}
	engine.food = Point{X: 5, Y: 4}
	engine.enemies = []enemySnake{
		{
			id:        1,
			direction: DirectionRight,
			snake: []Point{
				{X: 4, Y: 4},
				{X: 3, Y: 4},
				{X: 2, Y: 4},
			},
		},
	}

	beforeLength := len(engine.enemies[0].snake)
	engine.moveEnemies()
	state := engine.State()

	if len(state.Enemies) != 1 {
		t.Fatalf("enemy count = %d, want 1", len(state.Enemies))
	}

	if len(state.Enemies[0].Snake) != beforeLength+1 {
		t.Fatalf("enemy length = %d, want %d", len(state.Enemies[0].Snake), beforeLength+1)
	}
}

func TestPlayerEatingDropAddsScoreAndGrows(t *testing.T) {
	engine := NewEngine(Config{Width: 16, Height: 12}, rand.New(rand.NewSource(1)))
	engine.drops = []Point{{X: engine.players[0].snake[0].X + 1, Y: engine.players[0].snake[0].Y}}

	before := engine.State()
	engine.Turn(DirectionRight)
	engine.Tick()
	after := engine.State()

	if after.Score != 1 {
		t.Fatalf("score = %d, want 1", after.Score)
	}

	if len(after.Snake) != len(before.Snake)+1 {
		t.Fatalf("snake length = %d, want %d", len(after.Snake), len(before.Snake)+1)
	}

	if len(after.Drops) != 0 {
		t.Fatalf("drop count = %d, want 0", len(after.Drops))
	}
}
