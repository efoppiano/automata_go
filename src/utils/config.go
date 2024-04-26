package utils

type Config struct {
	CrosswalkProt         *Rectangle
	VehicleLaneProt       *Rectangle
	WaitingAreaProt       *Rectangle
	VehicleProt           *Rectangle
	StopLightCycle        int
	GreenLightTime        int
	PedestrianArrivalRate float64
	VehicleArrivalRate    float64
}

func NewConfig(
	crosswalkProt,
	vehicleLaneProt,
	waitingAreaProt,
	vehicleProt *Rectangle,
	stopLightCycle int,
	greenLightTime int,
	pedestrianArrivalRate,
	vehicleArrivalRate float64) *Config {
	return &Config{
		crosswalkProt,
		vehicleLaneProt,
		waitingAreaProt,
		vehicleProt,
		stopLightCycle,
		greenLightTime,
		pedestrianArrivalRate,
		vehicleArrivalRate,
	}
}

func (c *Config) TotalCols() int {
	return c.CrosswalkProt.Cols() + 2*c.WaitingAreaProt.Cols()
}

func (c *Config) TotalRows() int {
	return c.VehicleLaneProt.Rows()
}

func (c *Config) WalkingZoneProt() *Rectangle {
	return NewRectangle(c.CrosswalkProt.Rows(), c.TotalCols())
}

func (c *Config) Duplicate() *Config {
	return NewConfig(
		c.CrosswalkProt,
		c.VehicleLaneProt,
		c.WaitingAreaProt,
		c.VehicleProt,
		c.StopLightCycle,
		c.GreenLightTime,
		c.PedestrianArrivalRate,
		c.VehicleArrivalRate,
	)
}

func NewConfigFromEnv() *Config {
	crosswalkRows := GetEnvIntOrDefault("crosswalk_rows", 6)
	crosswalkCols := GetEnvIntOrDefault("crosswalk_cols", 42)
	waitingAreaCols := GetEnvIntOrDefault("waiting_area_cols", 1)
	vehicleLanes := GetEnvIntOrDefault("vehicle_lanes", 6)
	vehicleRows := GetEnvIntOrDefault("vehicle_rows", 6)
	vehicleCols := GetEnvIntOrDefault("vehicle_cols", 5)
	stopLightCycle := GetEnvIntOrDefault("stop_light_cycle", 90)
	greenLightTime := GetEnvIntOrDefault("green_light_time", 50)
	pedestrianArrivalRate := GetEnvFloatOrDefault("pedestrian_arrival_rate", 2000.0/(2*3600))
	vehicleArrivalRate := GetEnvFloatOrDefault("vehicle_arrival_rate", 1400.0/(6*3600))

	vehicleLaneCols := crosswalkCols / vehicleLanes
	crosswalkPrototype := NewRectangle(crosswalkRows, crosswalkCols)
	vehicleLanePrototype := NewRectangle(2*vehicleRows+crosswalkPrototype.Rows(), vehicleLaneCols)
	waitingAreaPrototype := NewRectangle(crosswalkRows, waitingAreaCols)
	vehiclePrototype := NewRectangle(vehicleRows, vehicleCols)

	return NewConfig(
		crosswalkPrototype,
		vehicleLanePrototype,
		waitingAreaPrototype,
		vehiclePrototype,
		stopLightCycle,
		greenLightTime,
		pedestrianArrivalRate,
		vehicleArrivalRate,
	)
}
