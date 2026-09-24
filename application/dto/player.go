package dto

import "time"

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
	Position     string  `json:"position"`
	JerseyNumber int     `json:"jersey_number"`
}

type PlayerPosition struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type GetListTeamPlayerItem struct {
	Id           string         `json:"id"`
	TeamId       string         `json:"team_id"`
	Name         string         `json:"name"`
	Position     PlayerPosition `json:"position"`
	JerseyNumber int            `json:"jersey_number"`
}

type ListPlayerItemTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetListPlayerItem struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Position     *PlayerPosition    `json:"position,omitempty"`
	JerseyNumber int                `json:"jersey_number"`
	Team         ListPlayerItemTeam `json:"team"`
}

type UpdatePlayerRequest struct {
	TeamId       *string  `json:"team_id"`
	Name         *string  `json:"name"`
	Height       *float64 `json:"height"`
	Weight       *float64 `json:"weight"`
	Position     *string  `json:"position"`
	JerseyNumber *int     `json:"jersey_number"`
}

type GetPlayerDetailResponse struct {
	Id           string              `json:"id"`
	TeamId       string              `json:"team_id"`
	Name         string              `json:"name"`
	Height       float64             `json:"height"`
	Weight       float64             `json:"weight"`
	Position     PlayerPosition      `json:"position"`
	JerseyNumber int                 `json:"jersey_number"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	IsDeleted    bool                `json:"is_deleted"`
	DeletedAt    time.Time           `json:"deleted_at,omitempty"`
	Team         *ListPlayerItemTeam `json:"team,omitempty"`
}
