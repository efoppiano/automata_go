package model

import (
	"go_automata/src/grid"
	"go_automata/src/utils"
)

/*
from typing import List

from rectangle import Rectangle, Point
from config import Config
from grid.grid import Grid

class Plotter:
    def __init__(self, grid: Grid, config: Config):
        self._config = config
        self._grid = grid
        self._bounds = Rectangle(self._config.crosswalk_prot.rows + 4*self._config.vehicle_prot.rows,
                                 self._config.total_cols)
        # self._bounds.move_down(self._config.vehicle_prot.rows)
        crosswalk_start_row = 3*self._config.vehicle_prot.rows
        crosswalk_start_col = self._config.waiting_area_prot.cols

        self._crosswalk_zone = self._config.crosswalk_prot.duplicate()
        self._crosswalk_zone.move_right(crosswalk_start_col)
        self._crosswalk_zone.move_down(crosswalk_start_row)

        self._waiting_area_zones: List[Rectangle] = []
        waiting_area_west_zone = self._config.waiting_area_prot.duplicate()
        waiting_area_west_zone.move_down(crosswalk_start_row)
        self._waiting_area_zones.append(waiting_area_west_zone)

        waiting_area_east_zone = self._config.waiting_area_prot.duplicate()
        waiting_area_east_zone.move_right(self._config.waiting_area_prot.cols + self._crosswalk_zone.cols)
        waiting_area_east_zone.move_down(crosswalk_start_row)
        self._waiting_area_zones.append(waiting_area_east_zone)

    def plot(self):
        self._grid.plot(self._plot_object, self._bounds)

    def _plot_object(self, point: Point, obj) -> str:
        _, col = point
        if obj is None:
             for waiting_area_zone in self._waiting_area_zones:
                if waiting_area_zone.is_inside(point):
                    return f"{'🔳'}"

             if self._crosswalk_zone.is_inside(point) and col % 4 <= 1:
                return f"{'⬜'}"
             else:
                return f"{'⬛'}"
        else:
            return obj._repr

*/

type Plotter struct {
	config           *utils.Config
	grid             *grid.Grid
	bounds           *utils.Rectangle
	crosswalkZone    *utils.Rectangle
	waitingAreaZones []*utils.Rectangle
}

func NewPlotter(grid *grid.Grid, config *utils.Config) *Plotter {
	crosswalkStartRow := config.VehicleProt.Rows()
	crosswalkStartCol := config.WaitingAreaProt.Cols()

	crosswalkZone := config.CrosswalkProt.Duplicate()
	crosswalkZone.MoveRight(crosswalkStartCol)
	crosswalkZone.MoveDown(crosswalkStartRow)

	waitingAreaZones := []*utils.Rectangle{}
	waitingAreaWestZone := config.WaitingAreaProt.Duplicate()
	waitingAreaWestZone.MoveDown(crosswalkStartRow)
	waitingAreaZones = append(waitingAreaZones, waitingAreaWestZone)

	waitingAreaEastZone := config.WaitingAreaProt.Duplicate()
	waitingAreaEastZone.MoveRight(config.WaitingAreaProt.Cols() + crosswalkZone.Cols())
	waitingAreaEastZone.MoveDown(crosswalkStartRow)
	waitingAreaZones = append(waitingAreaZones, waitingAreaEastZone)

	bounds := utils.NewRectangle(config.CrosswalkProt.Rows()+config.VehicleProt.Rows(), config.TotalCols())
	bounds.MoveDown(config.VehicleProt.Rows())

	return &Plotter{
		config:           config,
		grid:             grid,
		bounds:           bounds,
		crosswalkZone:    crosswalkZone,
		waitingAreaZones: waitingAreaZones,
	}
}

func (p *Plotter) Plot() {
	p.grid.Plot(p.plotObject, p.bounds)
}

func (p *Plotter) plotObject(point *utils.Point, obj interface{}) string {
	if obj == nil {
		for _, waitingAreaZone := range p.waitingAreaZones {
			if waitingAreaZone.IsInside(point) {
				return "🔳"
			}
		}

		if p.crosswalkZone.IsInside(point) && point.Y%4 <= 1 {
			return "⬜"
		} else {
			return "⬛"
		}
	} else {
		e := obj.(RoadEntity)
		return e.Repr()
	}
}
