package main

import (
	"fmt"
	"go_automata/src/generator"
	"go_automata/src/model"
	"go_automata/src/utils"
	"os"
	"os/exec"
	"time"
)

func CallClear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func linSpace(start, end float64, n int) []float64 {
	step := (end - start) / float64(n)
	result := make([]float64, 0)
	for i := 0; i < n; i++ {
		result = append(result, start+float64(i)*step)
	}
	return result
}

func buildConfigs() []*utils.Config {
	config := utils.NewFromEnvFile()
	configs := make([]*utils.Config, 0)
	for _, pedestrianArrivalRate := range linSpace(6000.0/(2*3600), 1000.0/(2*3600), 30) {
		for _, vehicleArrivalRate := range linSpace(1400.0/(6*3600), 200.0/(6*3600), 30) {
			configDup := config.Duplicate()
			configDup.PedestrianArrivalRate = pedestrianArrivalRate
			configDup.VehicleArrivalRate = vehicleArrivalRate
			configs = append(configs, configDup)
		}
	}
	return configs
}

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

func run(input_ch chan Input, results_ch chan *Result) {
	for input := range input_ch {
		i := input.i
		config := input.config

		fmt.Println("Starting automata...")
		start := time.Now()
		results := make([]int, 0)
		var j uint64
		for j = 0; j < 30; j++ {
			bbs := generator.NewBlumBlumShub(9000000 + i*100 + j)
			automata := model.NewAutomata(config, bbs)
			automata.AdvanceTo(3600)
			results = append(results, automata.Conflicts)
		}
		fmt.Println("Scenario ", i, " finished in", time.Since(start))

		pedestrianArrivalRate := config.PedestrianArrivalRate
		vehicleArrivalRate := config.VehicleArrivalRate
		total := 0
		for _, r := range results {
			total += r
		}
		average := float64(total) / float64(len(results))
		results_ch <- NewResult(pedestrianArrivalRate, vehicleArrivalRate, average)
	}
}

func main() {
	configs := buildConfigs()
	input_ch := make(chan Input, 10000)
	results_ch := make(chan *Result, 10000)
	var i uint64
	for i = 0; i < 20; i++ {
		go run(input_ch, results_ch)
	}

	for i, config := range configs {
		input_ch <- Input{uint64(i), config}
	}
	close(input_ch)
	t := time.Now()
	f, err := os.Create(fmt.Sprintf("results/%d-%d-%d-%d-%d-%d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("pedestrian_arrival_rate,vehicle_arrival_rate,conflicts\n")

	for i := 0; i < len(configs); i++ {
		r := <-results_ch
		f.WriteString(fmt.Sprintf("%d,%d,%f\n", int(r.PedestrianArrivalRate*2*3600), int(r.VehicleArrivalRate*6*3600), r.Conflicts))
	}
}
