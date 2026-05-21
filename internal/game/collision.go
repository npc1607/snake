package game

func (e *Engine) occupiedCells() map[Point]struct{} {
	occupied := make(map[Point]struct{}, len(e.snake)+len(e.drops))
	for _, segment := range e.snake {
		occupied[segment] = struct{}{}
	}

	for _, enemy := range e.enemies {
		for _, segment := range enemy.snake {
			occupied[segment] = struct{}{}
		}
	}

	for _, drop := range e.drops {
		occupied[drop] = struct{}{}
	}

	return occupied
}

func (e *Engine) occupiedWithoutEnemy(enemyIndex int) map[Point]struct{} {
	occupied := make(map[Point]struct{}, len(e.snake)+len(e.drops))
	for _, segment := range e.snake {
		occupied[segment] = struct{}{}
	}

	for index, enemy := range e.enemies {
		if index == enemyIndex {
			continue
		}

		for _, segment := range enemy.snake {
			occupied[segment] = struct{}{}
		}
	}

	for _, drop := range e.drops {
		occupied[drop] = struct{}{}
	}

	return occupied
}

func (e *Engine) hitsEnemy(point Point) bool {
	for _, enemy := range e.enemies {
		for _, segment := range enemy.snake {
			if segment == point {
				return true
			}
		}
	}

	return false
}

func (e *Engine) hitsEnemyAt(point Point, enemyIndex int, ignoreTail bool) bool {
	for index, enemy := range e.enemies {
		limit := len(enemy.snake)
		if index == enemyIndex && ignoreTail {
			limit--
		}

		for segmentIndex := 0; segmentIndex < limit; segmentIndex++ {
			if enemy.snake[segmentIndex] == point {
				return true
			}
		}
	}

	return false
}
