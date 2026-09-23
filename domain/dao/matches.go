package dao

import "time"

type Matches struct {
	Id         string
	MatchDate  *string
	MatchTime  *string
	HomeTeamId string
	AwayTeamId string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsDeleted  bool
	DeletedAt  time.Time
}
