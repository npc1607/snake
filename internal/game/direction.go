package game

// Direction describes where the snake moves on each tick.
type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

// IsOpposite reports whether two directions would reverse the snake into itself.
func (d Direction) IsOpposite(other Direction) bool {
	switch d {
	case DirectionUp:
		return other == DirectionDown
	case DirectionDown:
		return other == DirectionUp
	case DirectionLeft:
		return other == DirectionRight
	case DirectionRight:
		return other == DirectionLeft
	default:
		return false
	}
}
