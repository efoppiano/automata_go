package model

import "fmt"

type StopLightState int

const (
	Red StopLightState = iota
	Yellow
	Green
)

type StopLight struct {
	cycle           int
	greenLightTime  int
	yellowLightTime int
	timeToChange    int
	state           StopLightState
	prevState       StopLightState
}

func NewStopLight(cycle, greenLightTime int, yellowLightTime int) *StopLight {
	if greenLightTime+yellowLightTime >= cycle {
		panic("green light time must be less than cycle")
	}
	timeToChange := greenLightTime
	return &StopLight{cycle, greenLightTime, yellowLightTime, timeToChange, Green, Yellow}
}

func (sl *StopLight) Update() {
	sl.timeToChange--
	if sl.timeToChange == 0 {
		nextState := sl.getNextState()
		switch nextState {
		case Green:
			sl.timeToChange = sl.greenLightTime
		case Yellow:
			sl.timeToChange = sl.yellowLightTime
		case Red:
			sl.timeToChange = sl.cycle - sl.greenLightTime - sl.yellowLightTime
		}
		sl.prevState = sl.state
		sl.state = nextState
	}
}

func (sl *StopLight) getNextState() StopLightState {
	if sl.state == Green || sl.state == Red {
		return Yellow
	} else {
		if sl.PrevStateIsGreen() {
			return Red
		} else {
			return Green
		}
	}
}

func (sl *StopLight) IsGreen() bool {
	return sl.state == Green
}

func (sl *StopLight) IsYellow() bool {
	return sl.state == Yellow
}

func (sl *StopLight) IsRed() bool {
	return sl.state == Red
}

func (sl *StopLight) PrevStateIsGreen() bool {
	return sl.prevState == Green
}

func (sl *StopLight) PrevStateIsYellow() bool {
	return sl.prevState == Yellow
}

func (sl *StopLight) PrevStateIsRed() bool {
	return sl.prevState == Red
}

func (sl *StopLight) Show() {
	if sl.state == Green {
		fmt.Println("🟢", sl.timeToChange)
	} else if sl.state == Yellow {
		fmt.Println("🟡", sl.timeToChange)
	} else {
		fmt.Println("🔴", sl.timeToChange)
	}
}
