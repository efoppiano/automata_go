package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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

func getEnvIntOrDefault(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	r, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return r
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	r, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return r
}

func NewFromEnvFile() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	crosswalkRows := getEnvIntOrDefault("crosswalk_rows", 6)
	crosswalkCols := getEnvIntOrDefault("crosswalk_cols", 42)
	waitingAreaCols := getEnvIntOrDefault("waiting_area_cols", 1)
	vehicleLanes := getEnvIntOrDefault("vehicle_lanes", 6)
	vehicleRows := getEnvIntOrDefault("vehicle_rows", 6)
	vehicleCols := getEnvIntOrDefault("vehicle_cols", 5)
	stopLightCycle := getEnvIntOrDefault("stop_light_cycle", 90)
	greenLightTime := getEnvIntOrDefault("green_light_time", 50)
	pedestrianArrivalRate := getEnvFloatOrDefault("pedestrian_arrival_rate", 2000.0/(2*3600))
	vehicleArrivalRate := getEnvFloatOrDefault("vehicle_arrival_rate", 1400.0/(6*3600))

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
