package main

import (
	"fmt"
	"go_automata/src/utils"
)

type ScenarioConfig struct {
	InitialPedestrianArrivalRateHr int
	FinalPedestrianArrivalRateHr   int
	InitialVehicleArrivalRateHr    int
	FinalVehicleArrivalRateHr      int
	RunsPerSimulation              int
	SimulationTime                 int
}

func NewScenarioConfigFromEnv() *ScenarioConfig {
	initialPedestrianArrivalRateHr := utils.GetEnvIntOrDefault("INITIAL_PEDESTRIAN_ARRIVAL_RATE_HR", 1000)
	finalPedestrianArrivalRateHr := utils.GetEnvIntOrDefault("FINAL_PEDESTRIAN_ARRIVAL_RATE_HR", 6000)
	initialVehicleArrivalRateHr := utils.GetEnvIntOrDefault("INITIAL_VEHICLE_ARRIVAL_RATE_HR", 200)
	finalVehicleArrivalRateHr := utils.GetEnvIntOrDefault("FINAL_VEHICLE_ARRIVAL_RATE_HR", 1400)
	runsPerSimulation := utils.GetEnvIntOrDefault("RUNS_PER_SIMULATION", 30)
	simulationTime := utils.GetEnvIntOrDefault("SIMULATION_TIME", 3600)
	return &ScenarioConfig{
		initialPedestrianArrivalRateHr,
		finalPedestrianArrivalRateHr,
		initialVehicleArrivalRateHr,
		finalVehicleArrivalRateHr,
		runsPerSimulation,
		simulationTime,
	}
}

func (s *ScenarioConfig) Print() {
	println("Running with the following configuration:")
	println("Initial pedestrian arrival rate: ", s.InitialPedestrianArrivalRateHr, " cap/hr")
	println("Final pedestrian arrival rate: ", s.FinalPedestrianArrivalRateHr, " cap/hr")
	println("Initial vehicle arrival rate: ", s.InitialVehicleArrivalRateHr, " veh/hr")
	println("Final vehicle arrival rate: ", s.FinalVehicleArrivalRateHr, " veh/hr")
	println("Runs per simulation: ", s.RunsPerSimulation)
	println("Simulation time: ", s.SimulationTime, " seconds")
	println("Green light time: ", utils.GetEnvIntOrDefault("GREEN_LIGHT_TIME", 50), " seconds")
	fmt.Printf("Crosswalk width: %.1f meters\n", float64(utils.GetEnvIntOrDefault("CROSSWALK_ROWS", 6))/2)
}
