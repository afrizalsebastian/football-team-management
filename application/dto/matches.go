package dto

type CreateMatchRequest struct {
	HomeTeamId string `json:"home_team_id" validate:"required,nefield=AwayTeamId"`
	AwayTeamId string `json:"away_team_id" validate:"required,nefield=HomeTeamId"`
	Date       string `json:"date" validate:"required,ddmmyyyy"`
	Time       string `json:"time" validate:"required,hhmm"`
}

type CreateMatchResponse struct {
	Id         string `json:"id"`
	HomeTeamId string `json:"home_team_id"`
	AwayTeamId string `json:"away_team_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
}

type MatchTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type GetMatchResponse struct {
	Id         string              `json:"id"`
	HomeTeamId string              `json:"home_team_id"`
	AwayTeamId string              `json:"away_team_id"`
	Date       string              `json:"date"`
	Time       string              `json:"time"`
	Status     string              `json:"status"`
	HomeScore  int                 `json:"home_score"`
	AwayScore  int                 `json:"away_score"`
	HomeTeam   MatchTeam           `json:"home_team"`
	AwayTeam   MatchTeam           `json:"away_team"`
	Goals      []MatchGoalListItem `json:"goals"`
}

type RescheduleMatchRequest struct {
	Date string `json:"date" validate:"required,ddmmyyyy"`
	Time string `json:"time" validate:"required,hhmm"`
}

type MatchGoal struct {
	PlayerId   string `json:"player_id" validate:"required"`
	GoalMinute string `json:"goal_minute" validate:"required,goal"`
}

type MatchGoalListItem struct {
	Id         string            `json:"id"`
	GoalMinute string            `json:"goal_minute"`
	Player     GetListPlayerItem `json:"player"`
}

type MatchGoalsResponse struct {
	Goals    []MatchGoalListItem `json:"goals"`
	MatchId  string              `json:"match_id"`
	HomeTeam ListPlayerItemTeam  `json:"home_team"`
	AwayTeam ListPlayerItemTeam  `json:"away_team"`
}
