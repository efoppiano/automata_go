package utils

import (
	"os"
	"strconv"
)

func GetEnvIntOrDefault(key string, defaultValue int) int {
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

func GetEnvFloatOrDefault(key string, defaultValue float64) float64 {
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

func GetEnvStr(key string) *string {
	val := os.Getenv(key)
	if val == "" {
		return nil
	}
	return &val
}
