package grid

import "go_automata/src/utils"

type RelativeGrid struct {
	center utils.Point
	bounds *utils.Rectangle
	facing utils.Direction
	grid   *Grid
}

func NewRelativeGrid(center utils.Point, bounds *utils.Rectangle, facing utils.Direction, grid *Grid) *RelativeGrid {
	return &RelativeGrid{center, bounds, facing, grid}
}

func (rg *RelativeGrid) NewDisplaced(displacement *utils.RelativePosition) *RelativeGrid {
	newCenter := displacement.Apply(rg.facing, rg.center)
	return NewRelativeGrid(newCenter, rg.bounds, rg.facing, rg.grid)
}

func (rg *RelativeGrid) Facing() utils.Direction {
	return rg.facing
}

func (rg *RelativeGrid) Rows() int {
	return rg.bounds.Rows()
}

func (rg *RelativeGrid) Cols() int {
	return rg.bounds.Cols()
}

func (rg *RelativeGrid) IsIn(zone *utils.Rectangle) bool {
	return zone.IsInside(&rg.center)
}

func (rg *RelativeGrid) Fill(displacement *utils.RelativePosition, obj interface{}) {
	if !rg.IsInbounds(displacement) {
		panic("Attempted to fill an out of bounds cell")
	}

	point := displacement.Apply(rg.facing, rg.center)
	rg.grid.Fill(point.X, point.Y, obj)
}

func (rg *RelativeGrid) IsFill(displacement *utils.RelativePosition) bool {
	point := displacement.Apply(rg.facing, rg.center)
	return rg.grid.IsFill(point.X, point.Y)
}

func (rg *RelativeGrid) Get(displacement *utils.RelativePosition) interface{} {
	if !rg.IsInbounds(displacement) {
		return nil
	}
	point := displacement.Apply(rg.facing, rg.center)
	if !rg.grid.IsFill(point.X, point.Y) {
		return nil
	}
	return rg.grid.GetValue(point.X, point.Y)
}

func (rg *RelativeGrid) IsInbounds(displacement *utils.RelativePosition) bool {
	point := displacement.Apply(rg.facing, rg.center)
	if !rg.bounds.IsInside(&point) {
		return false
	}

	return 0 <= point.X && point.X < rg.grid.Rows() && 0 <= point.Y && point.Y < rg.grid.Cols()
}

func (rg *RelativeGrid) GetPrev(displacement *utils.RelativePosition, f func(interface{}) bool, maxChecks int) interface{} {
	if f == nil {
		f = func(interface{}) bool { return true }
	}

	point := displacement.Apply(rg.facing, rg.center)
	if rg.facing == utils.East {
		return rg.grid.GetPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.West {
		return rg.grid.GetNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.North {
		return rg.grid.GetVerticallyNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.South {
		return rg.grid.GetVerticallyPrev(point.X, point.Y, f, maxChecks)
	}
	return nil
}

func (rg *RelativeGrid) GetNext(displacement *utils.RelativePosition, f func(interface{}) bool, maxChecks int) interface{} {
	if f == nil {
		f = func(interface{}) bool { return true }
	}

	point := displacement.Apply(rg.facing, rg.center)
	if rg.facing == utils.East {
		return rg.grid.GetNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.West {
		return rg.grid.GetPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.North {
		return rg.grid.GetVerticallyPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.South {
		return rg.grid.GetVerticallyNext(point.X, point.Y, f, maxChecks)
	}
	return nil
}

func (rg *RelativeGrid) CalcDistToNext(displacement *utils.RelativePosition, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	point := displacement.Apply(rg.facing, rg.center)
	if rg.facing == utils.East {
		return rg.grid.CalcDistToNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.West {
		return rg.grid.CalcDistToPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.North {
		return rg.grid.CalcDistToVerticallyPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.South {
		return rg.grid.CalcDistToVerticallyNext(point.X, point.Y, f, maxChecks)
	}
	panic("Invalid direction")
}

func (rg *RelativeGrid) CalcDistToPrev(displacement *utils.RelativePosition, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	point := displacement.Apply(rg.facing, rg.center)
	if rg.facing == utils.East {
		return rg.grid.CalcDistToPrev(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.West {
		return rg.grid.CalcDistToNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.North {
		return rg.grid.CalcDistToVerticallyNext(point.X, point.Y, f, maxChecks)
	} else if rg.facing == utils.South {
		return rg.grid.CalcDistToVerticallyPrev(point.X, point.Y, f, maxChecks)
	}
	panic("Invalid direction")
}

func (rg *RelativeGrid) CalcDistToZone(displacement *utils.RelativePosition, zone utils.Rectangle) int {
	point := displacement.Apply(rg.facing, rg.center)
	return zone.DistanceTo(point)
}

func (rg *RelativeGrid) Clear(displacement *utils.RelativePosition) {
	point := displacement.Apply(rg.facing, rg.center)
	rg.grid.Clear(point.X, point.Y)
}

func (rg *RelativeGrid) Move(displacement *utils.RelativePosition) {
	if displacement.IsStill() {
		return
	}

	if !rg.IsFill(utils.NewRelativePosition(0, 0)) {
		panic("Attempted to move an empty cell")
	}

	if !rg.IsInbounds(displacement) {
		panic("Attempted to move out of bounds cell")
	}

	if rg.IsFill(displacement) {
		panic("Cell already fill")
	}

	me := rg.grid.GetValue(rg.center.X, rg.center.Y)
	point := displacement.Apply(rg.facing, rg.center)
	err := rg.grid.Fill(point.X, point.Y, me)
	if err != nil {
		displacement.Decrease()
		rg.Move(displacement)
		return
	}

	rg.grid.Clear(rg.center.X, rg.center.Y)
	rg.center = point
}
