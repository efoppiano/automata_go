package model

import "go_automata/src/utils"

type RoadEntity interface {
	IsCrossing() bool
	Think(crosswalkZone *utils.Rectangle, pedestrianStopLight *StopLight)
	Move(crosswalkZone *utils.Rectangle) bool
	IsVehicle() bool
	Repr() string
}
