package model

import (
	"go_automata/src/generator"
	"go_automata/src/grid"
	"go_automata/src/utils"
)

type Vehicle struct {
	repr             string
	vel              int
	crossing         bool
	desired_movement *utils.RelativePosition
	width            int
	length           int
	relative_origins []*grid.RelativeGrid
	driver_pos       *grid.RelativeGrid
	turning          bool
	generator        generator.Generator
}

func NewVehicle(origin *grid.RelativeGrid, prototype *utils.Rectangle, turning bool, generator generator.Generator) *Vehicle {
	repr_values := []string{"🟥", "🟧", "🟨", "🟩", "🟦", "🟪", "🟫"}
	i := generator.RandInt(0, len(repr_values))

	v := &Vehicle{
		vel:              10,
		crossing:         false,
		desired_movement: utils.Still(),
		width:            prototype.Cols(),
		length:           prototype.Rows(),
		turning:          turning,
		repr:             repr_values[i],
		generator:        generator,
	}
	v.buildGrids(origin)
	return v
}

func (v *Vehicle) buildGrids(origin *grid.RelativeGrid) {
	v.relative_origins = make([]*grid.RelativeGrid, 0)
	for i := 0; i < v.width; i++ {
		for j := 0; j < v.length; j++ {
			origin_ij := origin.NewDisplaced(utils.Right(i).Add(utils.Forward(j)))
			if i == 0 && j == v.length-1 {
				v.driver_pos = origin_ij
			} else {
				v.relative_origins = append(v.relative_origins, origin_ij)
			}
		}
	}
	v.driver_pos.Fill(utils.Still(), v)
	for _, relGridI := range v.relative_origins {
		part := NewVehiclePart(v, relGridI)
		relGridI.Fill(utils.Still(), part)
	}
}

func (v *Vehicle) Facing() utils.Direction {
	return v.driver_pos.Facing()
}

func (v *Vehicle) IsVehicle() bool {
	return true
}

func (v *Vehicle) IsCrossing() bool {
	return v.crossing
}

func (v *Vehicle) CanMove() bool {
	for i := 0; i < v.width; i++ {
		distToNext := v.driver_pos.CalcDistToNext(utils.Right(i), nil, v.vel)
		if distToNext != -1 {
			return false
		}
	}
	return true
}

func (v *Vehicle) IsEntityAhead() bool {
	for i := 0; i < v.width; i++ {
		entity := v.driver_pos.GetNext(utils.Right(i), nil, v.vel)
		if entity != nil {
			return true
		}
	}
	return false
}

func (v *Vehicle) IsPedestrianAhead() bool {
	for i := 0; i < v.width; i++ {
		entity := v.driver_pos.GetNext(utils.Right(i), nil, v.vel)
		switch entity.(type) {
		case *Pedestrian:
			return true
		}
	}
	return false
}

func (v *Vehicle) Think(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight) {
	if v.turning {
		v.thinkTurning(crosswalkZone, pedestrianStopLight)
	} else {
		v.thinkStraight(crosswalkZone, pedestrianStopLight)
	}
}

func (v *Vehicle) thinkStraight(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight) {
	if v.IsEntityAhead() || (pedestrianStopLight.IsGreen() && !v.crossing) {
		v.desired_movement = utils.Still()
	} else if v.crossing || pedestrianStopLight.IsRed() {
		v.desired_movement = utils.Forward(v.vel)
	} else {
		v.desired_movement = utils.Still()
	}
}

func (v *Vehicle) thinkTurning(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight) {
	if v.IsEntityAhead() {
		v.desired_movement = utils.Still()
	} else {
		v.desired_movement = utils.Forward(v.vel)
	}
}

func (v *Vehicle) Move(crosswalkZone *utils.Rectangle) bool {
	if v.desired_movement.IsStill() {
		return false
	}

	if v.IsPedestrianAhead() {
		for _, relGridI := range v.relative_origins {
			if relGridI.NewDisplaced(v.desired_movement).IsIn(crosswalkZone) {
				return true
			}
		}
		return false
	}

	if !v.driver_pos.IsInbounds(v.desired_movement) {
		v.Remove()
		return false
	}

	v.crossing = true
	v.driver_pos.Move(v.desired_movement)
	for _, relGridI := range v.relative_origins {
		relGridI.Move(v.desired_movement)
	}
	return false
}

func (v *Vehicle) Remove() {
	v.driver_pos.Clear(utils.Still())
	for _, relGridI := range v.relative_origins {
		relGridI.Clear(utils.Still())
	}
}

func (v *Vehicle) String() string {
	return v.repr
}

func (v *Vehicle) Repr() string {
	return v.repr
}

type VehiclePart struct {
	parent         *Vehicle
	relativeOrigin *grid.RelativeGrid
}

func NewVehiclePart(parent *Vehicle, relativeOrigin *grid.RelativeGrid) *VehiclePart {
	return &VehiclePart{
		parent:         parent,
		relativeOrigin: relativeOrigin,
	}
}

func (vp *VehiclePart) Facing() utils.Direction {
	return vp.relativeOrigin.Facing()
}

func (vp *VehiclePart) IsVehicle() bool {
	return true
}

func (vp *VehiclePart) IsCrossing() bool {
	return vp.parent.IsCrossing()
}

func (vp *VehiclePart) Think(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight) {
}

func (vp *VehiclePart) Move(crosswalkZone *utils.Rectangle) bool {
	return false
}

func (vp *VehiclePart) Repr() string {
	return vp.parent.Repr()
}
