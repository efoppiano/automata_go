package model

import (
	"go_automata/src/generator"
	"go_automata/src/grid"
	"go_automata/src/utils"
)

type VehicleLane struct {
	config          *utils.Config
	relGrid         *grid.RelativeGrid
	waitingVehicles int
	turning         bool
	generator       generator.Generator
}

func NewVehicleLane(config *utils.Config, relGrid *grid.RelativeGrid, turning bool, generator generator.Generator) *VehicleLane {
	return &VehicleLane{config, relGrid, 0, turning, generator}
}

func (vl *VehicleLane) generateVehicle() {
	newVehicles := vl.generator.Poi(vl.config.VehicleArrivalRate)
	vl.waitingVehicles += newVehicles
}

func (vl *VehicleLane) canPlaceVehicle() bool {
	for i := 0; i < vl.relGrid.Cols(); i++ {
		distToNext := vl.relGrid.CalcDistToNext(utils.Right(i), nil, vl.config.VehicleProt.Cols())
		if distToNext != -1 {
			return false
		}
	}
	return true
}

func (vl *VehicleLane) placeVehicle() {
	if vl.waitingVehicles == 0 || !vl.canPlaceVehicle() {
		return
	}

	offset := (vl.relGrid.Cols() - vl.config.VehicleProt.Cols()) / 2

	vehicleGrid := vl.relGrid.NewDisplaced(utils.Right(offset))
	NewVehicle(vehicleGrid, vl.config.VehicleProt, vl.turning, vl.generator)
	vl.waitingVehicles--
}

func (vl *VehicleLane) Update() {
	vl.generateVehicle()
	vl.placeVehicle()
}
