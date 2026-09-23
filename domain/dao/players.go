package dao

import (
	"time"
)

type Players struct {
	Id           string
	TeamId       string
	Name         *string
	Height       *float64
	Weight       *float64
	Position     *string
	JerseyNumber *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsDeleted    bool
	DeletedAt    time.Time
	Team         *Teams
}
