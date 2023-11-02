package cmd

import (
	"context"
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/ingest"
	"github.com/jaydeluca/otel-habits/pkg/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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

	ctx := context.Background()
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setting up meter provider")
	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	fmt.Println("Ingesting Habit data from Bear")
	habits := util.GenerateMetrics("habits", ingest.BearData(generationDayCount))

	fmt.Println("Ingesting Reading data from Google Sheets")
	reading := util.GenerateMetrics("reading", ingest.Sheets(generationDayCount))

	err = exp.Export(ctx, util.GenerateResourceMetrics([]metricdata.Metrics{
		*habits,
		*reading,
	}))
	if err != nil {
		fmt.Print(err)
	}

	return nil
}
