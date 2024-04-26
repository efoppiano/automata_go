package utils

type Direction int

const (
	East = iota
	West
	North
	South
)

func OppositeDirection(direction Direction) Direction {
	if direction == East {
		return West
	} else if direction == West {
		return East
	} else if direction == North {
		return South
	} else {
		return North
	}
}
