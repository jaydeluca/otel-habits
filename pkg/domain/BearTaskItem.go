package domain

import "time"

type BearTaskItem struct {
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
}
