package ingest

import "time"

type Timeseries struct {
	Date  time.Time
	Name  string
	Value int64
}
