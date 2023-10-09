package util

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/models"
	"math/rand"
	"time"
)

const (
	DateFormatString = "2006-01-02"
)

// RangeDate returns a date range function over start date to end date inclusive.
// After the end of the range, the range function returns a zero date,
// date.IsZero() is true.
func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

func GenerateDummyHabitData(recordsCount int) []models.Timeseries {
	fmt.Printf("-- Generating Dummy Habit Data --\n")

	goals := []string{
		"Workout",
		"Write Code",
		"Read Book",
	}

	entries := make([]models.Timeseries, 0)

	end := time.Now()
	start := end.AddDate(0, 0, -recordsCount)
	fmt.Println(
		"Generating data for time range: ", start.Format(DateFormatString), "-", end.Format(DateFormatString),
	)

	for rd := RangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		for _, goal := range goals {

			entry := models.Timeseries{
				Date:  date,
				Name:  goal,
				Value: int64(rand.Intn(2)),
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

func GenerateDummyReadingData(recordsCount int) []models.Timeseries {
	fmt.Printf("-- Generating Dummy Reading Data --\n")

	entries := make([]models.Timeseries, 0)

	end := time.Now()
	start := end.AddDate(0, 0, -recordsCount)
	fmt.Println(
		"Generating data for time range: ", start.Format(DateFormatString), "-", end.Format(DateFormatString),
	)

	for rd := RangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		entry := models.Timeseries{
			Date:  date,
			Name:  "reading",
			Value: int64(rand.Intn(60)),
		}
		entries = append(entries, entry)
	}

	return entries
}
