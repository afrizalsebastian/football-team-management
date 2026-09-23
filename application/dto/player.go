package dto

type CreatePlayerTeamRequest struct {
	Name         string  `json:"name" validate:"required"`
	Height       float64 `json:"height" validate:"required"`
	Weight       float64 `json:"weight" validate:"required"`
	Position     string  `json:"Position" validate:"required"`
	JerseyNumber int     `json:"jersey_number" validate:"required"`
}

type CreatePlayerTeamResponse struct {
	Id           string  `json:"id"`
	TeamId       string  `json:"team_id"`
	Name         string  `json:"name"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	Position     string  `json:"Position"`
	JerseyNumber int     `json:"jersey_number"`
}
