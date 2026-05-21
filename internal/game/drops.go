package game

func (e *Engine) dropIndex(point Point) int {
	for index, drop := range e.drops {
		if drop == point {
			return index
		}
	}

	return -1
}

func (e *Engine) removeDrop(index int) {
	e.drops = append(e.drops[:index], e.drops[index+1:]...)
}

func (e *Engine) scatterEnemy(enemy enemySnake) {
	for _, segment := range enemy.snake {
		if e.dropIndex(segment) >= 0 {
			continue
		}

		if segment == e.food {
			continue
		}

		e.drops = append(e.drops, segment)
	}
}
