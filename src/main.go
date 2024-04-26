package main

import (
	"fmt"
	"go_automata/src/utils"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
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

func buildConfigs(scenarioCfg *ScenarioConfig) []*utils.Config {
	baseCfg := utils.NewConfigFromEnv()
	configs := make([]*utils.Config, 0)
	initialPedestrianArrivalRate := float64(scenarioCfg.InitialPedestrianArrivalRateHr) / (2 * 3600)
	finalPedestrianArrivalRate := float64(scenarioCfg.FinalPedestrianArrivalRateHr) / (2 * 3600)
	initialVehicleArrivalRate := float64(scenarioCfg.InitialVehicleArrivalRateHr) / (6 * 3600)
	finalVehicleArrivalRate := float64(scenarioCfg.FinalVehicleArrivalRateHr) / (6 * 3600)

	for _, pedestrianArrivalRate := range linSpace(initialPedestrianArrivalRate, finalPedestrianArrivalRate, 30) {
		for _, vehicleArrivalRate := range linSpace(initialVehicleArrivalRate, finalVehicleArrivalRate, 30) {
			configDup := baseCfg.Duplicate()
			configDup.PedestrianArrivalRate = pedestrianArrivalRate
			configDup.VehicleArrivalRate = vehicleArrivalRate
			configs = append(configs, configDup)
		}
	}
	return configs
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	scenarioCfg := NewScenarioConfigFromEnv()
	scenarioCfg.Print()
	configs := buildConfigs(scenarioCfg)

	inputCh := make(chan Input, 10000)
	resultsCh := make(chan *Result, 10000)
	var i uint64
	for i = 0; i < 20; i++ {
		go run(scenarioCfg, inputCh, resultsCh)
	}

	for i, config := range configs {
		inputCh <- Input{uint64(i), config}
	}
	close(inputCh)

	t := time.Now()
	f, err := os.Create(fmt.Sprintf("results/%d-%d-%d-%d-%d-%d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("pedestrian_arrival_rate,vehicle_arrival_rate,conflicts\n")

	for i := 0; i < len(configs); i++ {
		r := <-resultsCh
		println("Received result ", i)
		f.WriteString(fmt.Sprintf("%d,%d,%f\n", int(r.PedestrianArrivalRate*2*3600), int(r.VehicleArrivalRate*6*3600), r.Conflicts))
	}
	close(resultsCh)
}
