package dao

import "time"

type AdminAccount struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
	IsDeleted bool
}
