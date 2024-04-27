package model

import (
	"go_automata/src/generator"
	"go_automata/src/grid"
	"go_automata/src/utils"
)

type Automata struct {
	Config              *utils.Config
	Grid                *grid.Grid
	CrosswalkZone       *utils.Rectangle
	Epoch               int
	Conflicts           int
	WaitingAreas        []*WaitingArea
	VehicleLanes        []*VehicleLane
	PedestrianStopLight *StopLight
	Plotter             *Plotter
	generator           generator.Generator
}

func NewAutomata(config *utils.Config, generator generator.Generator) *Automata {
	if config == nil {
		config = utils.NewConfigFromEnv()
	}

	totalRows := config.TotalRows()
	totalCols := config.TotalCols()
	grid := grid.NewGrid(totalRows, totalCols, generator)

	crosswalkZone := utils.NewRectangle(config.CrosswalkProt.Rows(), config.CrosswalkProt.Cols())
	crosswalkZone.MoveDown(config.VehicleProt.Rows())
	crosswalkZone.MoveRight(config.WaitingAreaProt.Cols())

	automata := &Automata{
		Config:              config,
		Grid:                grid,
		CrosswalkZone:       crosswalkZone,
		Epoch:               0,
		Conflicts:           0,
		PedestrianStopLight: NewStopLight(config.StopLightCycle, config.GreenLightTime, Green),
		Plotter:             NewPlotter(grid, config),
		generator:           generator,
	}

	automata.buildWaitingAreas()
	automata.buildVehicleLanes()

	return automata
}

func (a *Automata) buildWaitingAreas() {
	var walkingZone *utils.Rectangle
	if a.Config.WaitingAreaProt.Cols() > 0 {
		walkingZone = utils.NewRectangle(a.Config.CrosswalkProt.Rows(), a.Config.CrosswalkProt.Cols()+2)
		walkingZone.MoveRight(a.Config.WaitingAreaProt.Cols() - 1)
	} else {
		walkingZone = utils.NewRectangle(a.Config.CrosswalkProt.Rows(), a.Config.CrosswalkProt.Cols())
	}
	walkingZone.MoveDown(a.Config.VehicleProt.Rows())

	gridAreaWest := grid.NewRelativeGrid(walkingZone.UpperLeft, walkingZone, utils.East, a.Grid)
	gridAreaEast := grid.NewRelativeGrid(walkingZone.LowerRight, walkingZone, utils.West, a.Grid)

	a.WaitingAreas = []*WaitingArea{
		NewWaitingArea(a.Config.PedestrianArrivalRate, gridAreaWest, 100, a.generator),
		NewWaitingArea(a.Config.PedestrianArrivalRate, gridAreaEast, 100, a.generator),
	}
}

func (a *Automata) buildVehicleLanes() {
	vehicleLanesAmount := a.Config.CrosswalkProt.Cols() / a.Config.VehicleLaneProt.Cols()
	for i := 0; i < vehicleLanesAmount; i++ {
		vehicleLaneZone := a.Config.VehicleLaneProt.Duplicate()
		vehicleLaneZone.MoveRight(a.Config.WaitingAreaProt.Cols() + i*vehicleLaneZone.Cols())

		var facing utils.Direction
		var origin utils.Point
		if i < vehicleLanesAmount/2 {
			facing = utils.South
			origin = vehicleLaneZone.UpperRight()
		} else {
			facing = utils.North
			origin = vehicleLaneZone.LowerLeft()
		}

		grid := grid.NewRelativeGrid(origin, vehicleLaneZone, facing, a.Grid)

		var vehicleLane *VehicleLane
		if i == 0 || i == vehicleLanesAmount-1 {
			vehicleLane = NewVehicleLane(a.Config, grid, true, a.generator)
		} else {
			vehicleLane = NewVehicleLane(a.Config, grid, false, a.generator)
		}
		a.VehicleLanes = append(a.VehicleLanes, vehicleLane)
	}
}

func (a *Automata) Update() {
	a.PedestrianStopLight.Update()
	for _, waitingArea := range a.WaitingAreas {
		waitingArea.Update(a.PedestrianStopLight)
	}

	a.Grid.Apply(func(entity interface{}) {
		roadEntity := entity.(RoadEntity)
		roadEntity.Think(a.CrosswalkZone, a.PedestrianStopLight)
	})

	a.Grid.Apply(func(entity interface{}) {
		roadEntity := entity.(RoadEntity)
		conflictHappened := roadEntity.Move(a.CrosswalkZone)
		if conflictHappened {
			a.Conflicts++
		}
	})

	for _, vehicleLane := range a.VehicleLanes {
		vehicleLane.Update()
	}

	a.Epoch++
}

func (a *Automata) Show() {
	println("Epoch:", a.Epoch)
	println("Conflicts:", a.Conflicts)
	a.PedestrianStopLight.Show()
	println("Waiting at East:", a.WaitingAreas[0].waiting_pedestrians)
	println("Waiting at West:", a.WaitingAreas[1].waiting_pedestrians)
	a.Plotter.Plot()
}

func (a *Automata) AdvanceTo(epoch int) {
	for a.Epoch < epoch {
		a.Update()
	}
}
