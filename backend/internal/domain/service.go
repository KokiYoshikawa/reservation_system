package domain

import "time"

type Service struct {
	ID              int64
	Name            string
	DurationMinutes int
	Price           int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
