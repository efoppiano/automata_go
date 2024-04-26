package model

import (
	"go_automata/src/generator"
	"go_automata/src/grid"
	"go_automata/src/utils"
)

type WaitingArea struct {
	rel_grid            *grid.RelativeGrid
	waiting_pedestrians int
	arrival_rate        float64
	max_size            int
	generator           generator.Generator
}

func NewWaitingArea(arrival_rate float64, rel_grid *grid.RelativeGrid, max_size int, generator generator.Generator) *WaitingArea {
	return &WaitingArea{rel_grid, 0, arrival_rate, max_size, generator}
}

func (wa *WaitingArea) generatePedestrians() {
	if wa.waiting_pedestrians == wa.max_size {
		return
	}
	new_pedestrians := min(wa.max_size-wa.waiting_pedestrians, wa.generator.Poi(wa.arrival_rate))
	wa.waiting_pedestrians += new_pedestrians
}

func (wa *WaitingArea) canPlacePedestrian() bool {
	for i := 0; i < wa.rel_grid.Rows(); i++ {
		if !wa.rel_grid.IsFill(utils.Right(i)) {
			return true
		}
	}
	return false
}

func (wa *WaitingArea) placePedestrian() {
	rows := wa.rel_grid.Rows()

	possible_pos := wa.generator.RandInt(0, rows)
	for wa.rel_grid.IsFill(utils.Right(possible_pos)) {
		possible_pos = (possible_pos + 1) % rows
	}

	pedestrian_grid := wa.rel_grid.NewDisplaced(utils.Right(possible_pos))
	wa.rel_grid.Fill(utils.Right(possible_pos), NewPedestrian(pedestrian_grid, 0, "", wa.generator))
	wa.waiting_pedestrians--
}

func (wa *WaitingArea) placePedestrians() {
	for wa.waiting_pedestrians > 0 && wa.canPlacePedestrian() {
		wa.placePedestrian()
	}
}

func (wa *WaitingArea) Update(pedestrian_stop_light *StopLight) {
	wa.generatePedestrians()
	if pedestrian_stop_light == nil || pedestrian_stop_light.IsGreen() {
		wa.placePedestrians()
	}
}
