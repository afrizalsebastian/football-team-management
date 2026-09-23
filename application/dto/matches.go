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
