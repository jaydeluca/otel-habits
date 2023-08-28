package transform

import (
	"context"
	"github.com/jaydeluca/otel-habits/pkg/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"time"
)

var (
	res = resource.NewSchemaless(
		semconv.ServiceName("otel-habits"),
	)
)

func generateDataPoints(event domain.BearTaskItem) *metricdata.ResourceMetrics {
	mockData := metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "example", Version: "v0.0.1"},
				Metrics: []metricdata.Metrics{
					{
						Name:        "habits",
						Description: "Tracks counts of accomplishing a daily goal",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							IsMonotonic: true,
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("name", event.Name)),
									StartTime:  event.Date,
									Time:       time.Now(),
									Value:      1,
								},
							},
						},
					},
				},
			},
		},
	}
	return &mockData
}

func ConvertMetrics(events []domain.BearTaskItem) {
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
	for _, event := range events {
		_ = exp.Export(ctx, generateDataPoints(event))
	}
}
