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
	crosswalkRows := GetEnvIntOrDefault("CROSSWALK_ROWS", 6)
	crosswalkCols := GetEnvIntOrDefault("CROSSWALK_COLS", 42)
	waitingAreaCols := GetEnvIntOrDefault("WAITING_AREA_COLS", 1)
	vehicleLanes := GetEnvIntOrDefault("VEHICLE_LANES", 6)
	vehicleRows := GetEnvIntOrDefault("VEHICLE_ROWS", 6)
	vehicleCols := GetEnvIntOrDefault("VEHICLE_COLS", 5)
	stopLightCycle := GetEnvIntOrDefault("STOP_LIGHT_CYCLE", 90)
	greenLightTime := GetEnvIntOrDefault("GREEN_LIGHT_TIME", 50)
	pedestrianArrivalRate := GetEnvFloatOrDefault("PEDESTRIAN_ARRIVAL_RATE", 2000.0/(2*3600))
	vehicleArrivalRate := GetEnvFloatOrDefault("VEHICLE_ARRIVAL_RATE", 1400.0/(6*3600))

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
