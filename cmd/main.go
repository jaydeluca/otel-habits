package cmd

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/export"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
)

func Execute() error {
	fmt.Println("Starting App")
	events := ingest.Ingest()
	export.ExhaustMetrics(events)
	return nil
}
