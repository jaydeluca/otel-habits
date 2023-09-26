package cmd

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/export"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
)

func Execute() error {
	fmt.Println("Ingesting Habit data from Bear")
	export.ExhaustMetrics("habits", ingest.Ingest())

	fmt.Println("Ingesting Reading data from Google Sheets")
	export.ExhaustMetrics("reading", ingest.Sheets())
	return nil
}
