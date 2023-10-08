package export

import (
	"github.com/jaydeluca/otel-habits/pkg/models"
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
)

func generateDataPoints(metricName string, event models.Timeseries) *metricdata.ResourceMetrics {
	return &metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "habit-metrics", Version: "v0.0.1"},
				Metrics: []metricdata.Metrics{
					{
						Name:        metricName,
						Description: "",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							IsMonotonic: true,
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("name", event.Name)),
									StartTime:  event.Date,
									Time:       time.Now(),
									Value:      event.Value,
								},
							},
						},
					},
				},
			},
		},
	}
}
