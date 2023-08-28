package cmd

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
	"github.com/jaydeluca/otel-habits/pkg/transform"
)

func Execute() error {
	fmt.Println("Starting App")
	events := ingest.Ingest()
	transform.ConvertMetrics(events)
	return nil
}
