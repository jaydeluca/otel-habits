package cmd

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/export"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
)

func Execute() error {
	fmt.Println("Ingesting Habit models from Bear")
	export.ExhaustMetrics("habits", ingest.BearData())

	fmt.Println("Ingesting Reading models from Google Sheets")
	export.ExhaustMetrics("reading", ingest.Sheets())
	return nil
}
