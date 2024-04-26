package utils

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type Point struct {
	X, Y int
}

type Rectangle struct {
	UpperLeft, LowerRight Point
}

func NewRectangle(rows, cols int) *Rectangle {
	return &Rectangle{
		UpperLeft:  Point{0, 0},
		LowerRight: Point{rows - 1, cols - 1},
	}
}

func NewRectangleWithPoints(upperLeft, lowerRight Point) *Rectangle {
	return &Rectangle{
		UpperLeft:  upperLeft,
		LowerRight: lowerRight,
	}
}

func (r *Rectangle) IsInside(point *Point) bool {
	x, y := point.X, point.Y
	x1, y1 := r.UpperLeft.X, r.UpperLeft.Y
	x2, y2 := r.LowerRight.X, r.LowerRight.Y
	return x1 <= x && x <= x2 && y1 <= y && y <= y2
}

func (r *Rectangle) StartCol() int {
	return r.UpperLeft.Y
}

func (r *Rectangle) StartRow() int {
	return r.UpperLeft.X
}

func (r *Rectangle) EndCol() int {
	return r.LowerRight.Y
}

func (r *Rectangle) EndRow() int {
	return r.LowerRight.X
}

func (r *Rectangle) Cols() int {
	return r.EndCol() - r.StartCol() + 1
}

func (r *Rectangle) Rows() int {
	return r.EndRow() - r.StartRow() + 1
}

func (r *Rectangle) UpperRight() Point {
	return Point{r.UpperLeft.X, r.LowerRight.Y}
}

func (r *Rectangle) LowerLeft() Point {
	return Point{r.LowerRight.X, r.UpperLeft.Y}
}

func (r *Rectangle) MoveUp(rows int) {
	r.UpperLeft.X -= rows
	r.LowerRight.X -= rows
}

func (r *Rectangle) MoveDown(rows int) {
	r.MoveUp(-rows)
}

func (r *Rectangle) MoveLeft(cols int) {
	r.UpperLeft.Y -= cols
	r.LowerRight.Y -= cols
}

func (r *Rectangle) MoveRight(cols int) {
	r.MoveLeft(-cols)
}

func (r *Rectangle) Duplicate() *Rectangle {
	return NewRectangleWithPoints(r.UpperLeft, r.LowerRight)
}

func (r *Rectangle) DistanceTo(point Point) int {
	row, col := point.X, point.Y

	if r.StartRow() <= row && row <= r.EndRow() {
		return int(min(abs(r.StartCol()-col), abs(r.EndCol()-col))) - 1
	}
	if r.StartCol() <= col && col <= r.EndCol() {
		return int(min(abs(r.StartRow()-row), abs(r.EndRow()-row))) - 1
	}
	panic("Cannot calculate distance to point from rectangle")
}
