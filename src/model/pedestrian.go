package model

import (
	"go_automata/src/generator"
	"go_automata/src/grid"
	"go_automata/src/utils"
)

type Pedestrian struct {
	desired_displacement *utils.RelativePosition
	rel_grid             *grid.RelativeGrid
	crossing             bool
	vel                  int
	repr                 string
	generator            generator.Generator
}

func NewPedestrian(rel_grid *grid.RelativeGrid, velocity int, repr string, generator generator.Generator) *Pedestrian {
	p := &Pedestrian{rel_grid: rel_grid, crossing: false, generator: generator}
	p.desired_displacement = utils.Still()

	if velocity != 0 {
		p.vel = velocity
	} else {
		p.vel = p.generateVelocity()
	}

	if repr != "" {
		p.repr = repr
	} else {
		possible_values := []string{"😀", "😁", "🙃", "🤔", "😶", "🙄", "😎"}
		i := p.generator.RandInt(0, len(possible_values))
		p.repr = possible_values[i]
	}

	return p
}

func (p *Pedestrian) Facing() utils.Direction {
	return p.rel_grid.Facing()
}

func (p *Pedestrian) IsVehicle() bool {
	return false
}

func (p *Pedestrian) IsCrossing() bool {
	return p.crossing
}

func (p *Pedestrian) String() string {
	return p.repr
}

func (p *Pedestrian) generateVelocity() int {
	n := p.generator.Random()

	if n > 0.978 {
		return 6
	} else if n > 0.93 {
		return 5
	} else if n > 0.793 {
		return 4
	} else if n > 0.273 {
		return 3
	} else {
		return 2
	}
}

func (p *Pedestrian) CanMoveForward() bool {
	if !p.rel_grid.IsInbounds(utils.Forward(1)) {
		return true
	}

	dist_to_next := p.rel_grid.CalcDistToNext(utils.Still(), func(ent interface{}) bool {
		switch pedestrian := ent.(type) {
		case *Pedestrian:
			return pedestrian.IsCrossing() && pedestrian.Facing() == p.Facing()
		default:
			return false
		}
	}, 1)

	return dist_to_next == -1
}

func (p *Pedestrian) CanDoLateralMovement(toRight bool) bool {
	displacement := utils.Right(1)
	if !toRight {
		displacement = utils.Left(1)
	}

	if !p.rel_grid.IsFill(utils.Forward(1)) {
		return false
	}

	if !p.rel_grid.IsInbounds(displacement) {
		return false
	}

	if p.rel_grid.IsFill(displacement) {
		return false
	}

	dist := p.rel_grid.CalcDistToNext(displacement, func(ent interface{}) bool {
		switch pedestrian := ent.(type) {
		case *Pedestrian:
			return pedestrian.IsCrossing() && pedestrian.Facing() == utils.OppositeDirection(p.Facing())
		default:
			return false
		}

	}, p.vel)

	if dist != -1 {
		return false
	}

	dist_to_prev := p.rel_grid.CalcDistToPrev(displacement, func(ent interface{}) bool {
		switch pedestrian := ent.(type) {
		case *Pedestrian:
			return pedestrian.IsCrossing() && pedestrian.Facing() == p.Facing()
		default:
			return false
		}
	}, 6)

	if dist_to_prev == -1 {
		return true
	}

	prev := p.rel_grid.GetPrev(displacement, func(ent interface{}) bool {
		switch pedestrian := ent.(type) {
		case *Pedestrian:
			return pedestrian.IsCrossing() && pedestrian.Facing() == p.Facing()
		default:
			return false
		}
	}, 6)

	return prev.(*Pedestrian).vel < p.vel
}

func (p *Pedestrian) CanMoveLeft() bool {
	return p.CanDoLateralMovement(false)
}

func (p *Pedestrian) CanMoveRight() bool {
	return p.CanDoLateralMovement(true)
}

func (p *Pedestrian) GetPosForward() *utils.RelativePosition {
	dist_to_next := p.rel_grid.CalcDistToNext(utils.Still(), func(ent interface{}) bool {
		switch pedestrian := ent.(type) {
		case *Pedestrian:
			return pedestrian.IsCrossing() && pedestrian.Facing() == p.Facing()
		default:
			return false
		}
	}, 0)

	if dist_to_next == -1 || dist_to_next > p.vel {
		return utils.Forward(p.vel)
	}

	return utils.Forward(dist_to_next)
}

func (p *Pedestrian) GetPosLeftRightRandom() *utils.RelativePosition {
	n := p.generator.Random()
	if n > 0.5 {
		return utils.Left(1)
	} else {
		return utils.Right(1)
	}
}

func (p *Pedestrian) Think(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight) {
	if pedestrianStopLight.IsYellow() && pedestrianStopLight.PrevStateIsGreen() {
		if !p.crossing {
			p.desired_displacement = utils.Still()
			return
		} else {
			p.vel = 6
			p.repr = "😰"
		}
	}

	if pedestrianStopLight.IsRed() {
		if !p.rel_grid.IsIn(crosswalkZone) {
			p.desired_displacement = utils.Still()
			return
		} else {
			p.vel = 6
			p.repr = "😰"
		}
	}

	if p.CanMoveForward() {
		p.desired_displacement = p.GetPosForward()
	} else {
		can_move_left := p.CanMoveLeft()
		can_move_right := p.CanMoveRight()

		if can_move_left && !can_move_right {
			p.desired_displacement = utils.Left(1)
		} else if can_move_right && !can_move_left {
			p.desired_displacement = utils.Right(1)
		} else if can_move_left && can_move_right {
			p.desired_displacement = p.GetPosLeftRightRandom()
		} else {
			p.desired_displacement = utils.Still()
		}
	}
}

func (p *Pedestrian) Move(crosswalkZone *utils.Rectangle) bool {
	if !p.rel_grid.IsInbounds(p.desired_displacement) {
		p.rel_grid.Clear(utils.Still())
		return false
	}

	if p.desired_displacement.IsStill() {
		return false
	}

	p.crossing = true
	if !p.rel_grid.NewDisplaced(p.desired_displacement).IsIn(crosswalkZone) {
		p.rel_grid.Clear(utils.Still())
		return false
	}

	for !p.desired_displacement.IsStill() && p.rel_grid.IsFill(p.desired_displacement) {
		p.desired_displacement.Decrease()
	}

	if p.desired_displacement.IsStill() {
		return false
	}

	p.rel_grid.Move(p.desired_displacement)
	return false
}

func (p *Pedestrian) Repr() string {
	return p.repr
}
