package dao

import "time"

type Goals struct {
	Id         string
	MatchId    string
	PlayerId   string
	GoalMinute string
	CreatedAt  time.Time
	IsDeleted  bool
	DeletedAt  time.Time

	Match  *Matches
	Player *Players
}
