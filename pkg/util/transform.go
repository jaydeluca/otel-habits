package util

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"time"
)

var (
	res = resource.NewSchemaless(
		semconv.ServiceName("otel-habits"),
	)
	scope = instrumentation.Scope{Name: "goals", Version: "v0.0.1"}
)

func GenerateDataPoint(date time.Time, name string, value float64) *metricdata.DataPoint[float64] {
	return &metricdata.DataPoint[float64]{
		Attributes: attribute.NewSet(attribute.String("name", name)),
		StartTime:  date,
		Time:       time.Now(),
		Value:      value,
	}
}

func GenerateMetrics(metricName string, dataPoints []metricdata.DataPoint[float64]) *metricdata.Metrics {
	return &metricdata.Metrics{
		Name:        metricName,
		Description: "",
		Unit:        "1",
		Data: metricdata.Sum[float64]{
			IsMonotonic: true,
			Temporality: metricdata.DeltaTemporality,
			DataPoints:  dataPoints,
		},
	}
}

func GenerateResourceMetrics(metrics []metricdata.Metrics) *metricdata.ResourceMetrics {
	return &metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope:   scope,
				Metrics: metrics,
			},
		},
	}
}
