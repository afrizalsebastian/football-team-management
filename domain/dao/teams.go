package dao

import "time"

type Teams struct {
	Id          string
	Name        *string
	Logo        *string
	FoundedYear *string
	Address     *string
	City        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsDeleted   bool
	DeletedAt   time.Time
}
