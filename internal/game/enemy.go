package game

const (
	defaultMaxSnakes               = 6
	defaultEnemySpawnIntervalTicks = 316
	enemyStartLength               = 4
)

// EnemyState is a read-only enemy snake snapshot.
type EnemyState struct {
	ID        int
	Snake     []Point
	Direction Direction
}

type enemySnake struct {
	id        int
	snake     []Point
	direction Direction
}

func enemyStates(enemies []enemySnake) []EnemyState {
	states := make([]EnemyState, len(enemies))
	for index, enemy := range enemies {
		snake := make([]Point, len(enemy.snake))
		copy(snake, enemy.snake)

		states[index] = EnemyState{
			ID:        enemy.id,
			Snake:     snake,
			Direction: enemy.direction,
		}
	}

	return states
}

func playerStates(players []playerSnake) []PlayerState {
	states := make([]PlayerState, len(players))
	for index, player := range players {
		snake := make([]Point, len(player.snake))
		copy(snake, player.snake)

		states[index] = PlayerState{
			ID:        player.id,
			Snake:     snake,
			Direction: player.direction,
			Alive:     player.alive,
		}
	}

	return states
}
