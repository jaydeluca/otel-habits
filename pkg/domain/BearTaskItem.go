package domain

import "time"

type BearTaskItem struct {
	Date      time.Time
	Name      string
	Completed bool
}
