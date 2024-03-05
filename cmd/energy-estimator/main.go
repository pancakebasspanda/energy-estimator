package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"energy-estimator/processor"
	"energy-estimator/storage"
)

func main() {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app_name", "energy-estimator").
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	store := storage.New()
	processor := processor.New(log, store)

	processor.Process()

	// estimate the power usage
	fmt.Printf("Estimated energy used: %v Wh", store.CalculateEstimatedPowerConsumption())
}
