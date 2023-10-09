package cmd

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/export"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
	"os"
	"strconv"
)

func Execute() error {
	generationDayCount := 0
	var err error
	dayCountVar := os.Getenv("GENERATE_DATA_DAY_COUNT")
	if dayCountVar != "" {
		generationDayCount, err = strconv.Atoi(dayCountVar)
		if err != nil {
			panic(err)
		}
		fmt.Sprintf("Env var set to generate %d days worth of data", generationDayCount)
	}

	fmt.Println("Ingesting Habit data from Bear")
	export.ExhaustMetrics("habits", ingest.BearData(generationDayCount))

	fmt.Println("Ingesting Reading data from Google Sheets")
	export.ExhaustMetrics("reading", ingest.Sheets(generationDayCount))
	return nil
}
