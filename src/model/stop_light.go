package model

import "fmt"

type StopLightState int

const (
	Red StopLightState = iota
	Green
)

type StopLight struct {
	cycle          int
	greenLightTime int
	timeToChange   int
	state          StopLightState
}

func NewStopLight(cycle, greenLightTime int, initialState StopLightState) *StopLight {
	if greenLightTime >= cycle {
		panic("green light time must be less than cycle")
	}
	timeToChange := greenLightTime
	if initialState == Red {
		timeToChange = cycle - greenLightTime
	}
	return &StopLight{cycle, greenLightTime, timeToChange, initialState}
}

func (sl *StopLight) Update() {
	sl.timeToChange--
	if sl.timeToChange == 0 {
		if sl.state == Red {
			sl.state = Green
			sl.timeToChange = sl.greenLightTime
		} else {
			sl.state = Red
			sl.timeToChange = sl.cycle - sl.greenLightTime
		}
	}
}

func (sl *StopLight) IsGreen() bool {
	return sl.state == Green
}

func (sl *StopLight) IsRed() bool {
	return sl.state == Red
}

func (sl *StopLight) Show() {
	if sl.state == Green {
		fmt.Println("🟢", sl.timeToChange)
	} else {
		fmt.Println("🔴", sl.timeToChange)
	}
}
