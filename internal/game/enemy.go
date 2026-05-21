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
