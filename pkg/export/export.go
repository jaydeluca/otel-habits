package export

import (
	"context"
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
)

func ExhaustMetrics(events []domain.BearTaskItem) {
	ctx := context.Background()
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	// export to collector
	fmt.Print("Exporting metrics")
	for _, event := range events {
		_ = exp.Export(ctx, generateDataPoints(event))
	}
}
