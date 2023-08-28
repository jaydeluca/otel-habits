package util

import (
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/domain"
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

// For generating random booleans in dummy data script
type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

func NewBoolean() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

func GenerateDummyData(recordsCount int) []domain.BearTaskItem {
	fmt.Printf("-- Generating Dummy Data --\n")

	goals := []string{
		"Workout",
		"Write Code",
		"Walk Dogs",
	}

	entries := make([]domain.BearTaskItem, 0)

	end := time.Now()
	start := end.AddDate(0, 0, -recordsCount)
	fmt.Println(
		"Generating data for time range: ", start.Format(DateFormatString), "-", end.Format(DateFormatString),
	)

	r := NewBoolean()

	for rd := RangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		for _, goal := range goals {

			entry := domain.BearTaskItem{
				Date:      date,
				Name:      goal,
				Completed: r.Bool(),
			}
			entries = append(entries, entry)
		}
	}

	return entries
}
