package main

import (
	"fmt"
	"go_automata/src/generator"
	"go_automata/src/model"
	"go_automata/src/utils"
	"time"
)

type Result struct {
	PedestrianArrivalRate float64
	VehicleArrivalRate    float64
	Conflicts             float64
}

type Input struct {
	i      uint64
	config *utils.Config
}

func NewResult(pedestrianArrivalRate, vehicleArrivalRate float64, conflicts float64) *Result {
	return &Result{
		PedestrianArrivalRate: pedestrianArrivalRate,
		VehicleArrivalRate:    vehicleArrivalRate,
		Conflicts:             conflicts,
	}
}

func run(cfg *ScenarioConfig, inputCh chan Input, resultsCh chan *Result) {
	for input := range inputCh {
		i := input.i
		config := input.config

		println("Starting automata...")
		start := time.Now()
		results := make([]int, 0)
		var j uint64
		for j = 0; j < uint64(cfg.RunsPerSimulation); j++ {
			bbs := generator.NewBlumBlumShub(9000000 + i*100 + j)
			automata := model.NewAutomata(config, bbs)
			automata.AdvanceTo(cfg.SimulationTime)
			results = append(results, automata.Conflicts)
		}
		fmt.Println("Scenario", i, "finished in", time.Since(start))

		pedestrianArrivalRate := config.PedestrianArrivalRate
		vehicleArrivalRate := config.VehicleArrivalRate
		total := 0
		for _, r := range results {
			total += r
		}
		average := float64(total) / float64(len(results))
		fmt.Printf("Pedestrian arrival rate: %.2f, Vehicle arrival rate: %.2f, Average conflicts: %.2f\n", pedestrianArrivalRate, vehicleArrivalRate, average)
		resultsCh <- NewResult(pedestrianArrivalRate, vehicleArrivalRate, average)
	}
}
