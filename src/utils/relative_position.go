package utils

import "fmt"

type RelativePosition struct {
	forward int
	right   int
}

func NewRelativePosition(forward, right int) *RelativePosition {
	return &RelativePosition{forward, right}
}

func (rp *RelativePosition) IsStill() bool {
	return rp.forward == 0 && rp.right == 0
}

func (rp *RelativePosition) Decrease() {
	if rp.forward > 0 {
		rp.forward--
	} else if rp.forward < 0 {
		rp.forward++
	}
	if rp.right > 0 {
		rp.right--
	} else if rp.right < 0 {
		rp.right++
	}
}

func (rp *RelativePosition) Add(other *RelativePosition) *RelativePosition {
	return &RelativePosition{rp.forward + other.forward, rp.right + other.right}
}

func (rp *RelativePosition) Apply(facing Direction, center Point) Point {
	switch facing {
	case East:
		return Point{center.X + rp.right, center.Y + rp.forward}
	case West:
		return Point{center.X - rp.right, center.Y - rp.forward}
	case North:
		return Point{center.X - rp.forward, center.Y + rp.right}
	case South:
		return Point{center.X + rp.forward, center.Y - rp.right}
	default:
		panic(fmt.Sprintf("Invalid direction %v", facing))
	}
}

func Forward(amount int) *RelativePosition {
	return &RelativePosition{amount, 0}
}

func Backward(amount int) *RelativePosition {
	return &RelativePosition{-amount, 0}
}

func Right(amount int) *RelativePosition {
	return &RelativePosition{0, amount}
}

func Left(amount int) *RelativePosition {
	return &RelativePosition{0, -amount}
}

func Still() *RelativePosition {
	return &RelativePosition{0, 0}
}
