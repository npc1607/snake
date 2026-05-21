package game

func shortestDirection(start Point, target Point, blocked map[Point]struct{}, width int, height int) (Direction, bool) {
	if start == target {
		return DirectionRight, false
	}

	queue := []Point{start}
	seen := map[Point]struct{}{
		start: {},
	}
	previous := make(map[Point]Point)
	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, direction := range orderedDirections() {
			next := current.move(direction)
			if next.X < 0 || next.X >= width || next.Y < 0 || next.Y >= height {
				continue
			}

			if _, ok := seen[next]; ok {
				continue
			}

			if _, ok := blocked[next]; ok && next != target {
				continue
			}

			seen[next] = struct{}{}
			previous[next] = current
			if next == target {
				found = true
				queue = nil
				break
			}

			queue = append(queue, next)
		}
	}

	if !found {
		return DirectionRight, false
	}

	step := target
	for previous[step] != start {
		step = previous[step]
	}

	return directionBetween(start, step), true
}

func orderedDirections() []Direction {
	return []Direction{
		DirectionUp,
		DirectionLeft,
		DirectionRight,
		DirectionDown,
	}
}

func directionBetween(from Point, to Point) Direction {
	switch {
	case to.X < from.X:
		return DirectionLeft
	case to.X > from.X:
		return DirectionRight
	case to.Y < from.Y:
		return DirectionUp
	default:
		return DirectionDown
	}
}
