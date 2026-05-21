package game

func (e *Engine) maybeSpawnEnemy() {
	e.ticksSinceEnemy++
	if e.ticksSinceEnemy < e.enemySpawnIntervalTicks {
		return
	}

	e.ticksSinceEnemy = 0
	if len(e.enemies)+1 >= e.maxSnakes {
		return
	}

	e.spawnEnemy()
}

func (e *Engine) spawnEnemy() {
	occupied := e.occupiedCells()
	for attempt := 0; attempt < 80; attempt++ {
		head := Point{
			X: e.random.Intn(e.width),
			Y: e.random.Intn(e.height),
		}

		direction := DirectionLeft
		snake := buildEnemySnake(head, direction)
		if !e.canPlaceSnake(snake, occupied) {
			continue
		}

		e.nextEnemyID++
		e.enemies = append(e.enemies, enemySnake{
			id:        e.nextEnemyID,
			snake:     snake,
			direction: direction,
		})
		return
	}
}

func buildEnemySnake(head Point, direction Direction) []Point {
	snake := make([]Point, 0, enemyStartLength)
	snake = append(snake, head)

	opposite := oppositeDirection(direction)
	current := head
	for len(snake) < enemyStartLength {
		current = current.move(opposite)
		snake = append(snake, current)
	}

	return snake
}

func (e *Engine) canPlaceSnake(snake []Point, occupied map[Point]struct{}) bool {
	for _, segment := range snake {
		if e.hitsWall(segment) {
			return false
		}

		if _, ok := occupied[segment]; ok {
			return false
		}
	}

	return true
}

func (e *Engine) moveEnemies() {
	if len(e.enemies) == 0 {
		return
	}

	nextEnemies := make([]enemySnake, 0, len(e.enemies))
	for index := range e.enemies {
		enemy := e.enemies[index]
		direction := e.enemyDirection(index)
		nextHead := enemy.snake[0].move(direction)

		if nextHead == e.snake[0] {
			e.status = StatusGameOver
			e.recordScore()
			return
		}

		if e.enemyCrashes(nextHead, index) {
			e.scatterEnemy(enemy)
			continue
		}

		if nextHead == e.food {
			e.food = e.nextFood()
		}

		nextSnake := make([]Point, 0, len(enemy.snake))
		nextSnake = append(nextSnake, nextHead)
		nextSnake = append(nextSnake, enemy.snake[:len(enemy.snake)-1]...)
		enemy.direction = direction
		enemy.snake = nextSnake
		nextEnemies = append(nextEnemies, enemy)
	}

	e.enemies = nextEnemies
}

func (e *Engine) enemyDirection(enemyIndex int) Direction {
	enemy := e.enemies[enemyIndex]
	head := enemy.snake[0]
	blocked := e.occupiedWithoutEnemy(enemyIndex)

	target := e.enemyTarget(enemyIndex)
	if direction, ok := shortestDirection(head, target, blocked, e.width, e.height); ok {
		if !enemy.direction.IsOpposite(direction) {
			return direction
		}
	}

	return e.safeEnemyDirection(enemyIndex)
}

func (e *Engine) enemyTarget(enemyIndex int) Point {
	enemyHead := e.enemies[enemyIndex].snake[0]
	playerHead := e.snake[0]
	foodDistance := manhattan(enemyHead, e.food)
	playerDistance := manhattan(enemyHead, playerHead)

	if playerDistance <= foodDistance+3 {
		return playerHead
	}

	return e.food
}

func (e *Engine) safeEnemyDirection(enemyIndex int) Direction {
	enemy := e.enemies[enemyIndex]
	directions := append([]Direction{enemy.direction}, orderedDirections()...)
	for _, direction := range directions {
		if enemy.direction.IsOpposite(direction) {
			continue
		}

		nextHead := enemy.snake[0].move(direction)
		if !e.enemyCrashes(nextHead, enemyIndex) {
			return direction
		}
	}

	return enemy.direction
}

func (e *Engine) enemyCrashes(point Point, enemyIndex int) bool {
	if e.hitsWall(point) {
		return true
	}

	if e.hitsPlayerBody(point) {
		return true
	}

	return e.hitsEnemyAt(point, enemyIndex, true)
}

func (e *Engine) hitsPlayerBody(point Point) bool {
	for _, segment := range e.snake[1:] {
		if segment == point {
			return true
		}
	}

	return false
}

func (e *Engine) nextEnemyCountdown() int {
	if len(e.enemies)+1 >= e.maxSnakes {
		return 0
	}

	remaining := e.enemySpawnIntervalTicks - e.ticksSinceEnemy
	if remaining < 0 {
		return 0
	}

	return remaining
}

func manhattan(a Point, b Point) int {
	distanceX := a.X - b.X
	if distanceX < 0 {
		distanceX = -distanceX
	}

	distanceY := a.Y - b.Y
	if distanceY < 0 {
		distanceY = -distanceY
	}

	return distanceX + distanceY
}

func oppositeDirection(direction Direction) Direction {
	switch direction {
	case DirectionUp:
		return DirectionDown
	case DirectionDown:
		return DirectionUp
	case DirectionLeft:
		return DirectionRight
	case DirectionRight:
		return DirectionLeft
	default:
		return DirectionRight
	}
}
