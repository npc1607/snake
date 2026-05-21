package game

// Point identifies one board cell.
type Point struct {
	X int
	Y int
}

func (p Point) move(direction Direction) Point {
	switch direction {
	case DirectionUp:
		return Point{X: p.X, Y: p.Y - 1}
	case DirectionDown:
		return Point{X: p.X, Y: p.Y + 1}
	case DirectionLeft:
		return Point{X: p.X - 1, Y: p.Y}
	case DirectionRight:
		return Point{X: p.X + 1, Y: p.Y}
	default:
		return p
	}
}
