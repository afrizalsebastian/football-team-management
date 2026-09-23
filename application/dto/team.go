package dto

import "time"

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Logo        string `json:"logo" validate:"required"`
	FoundedYear string `json:"founded_year" validate:"required"`
	Address     string `json:"address" validate:"required"`
	City        string `json:"city" validate:"required"`
}

type CreateTeamResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	FoundedYear string `json:"founded_year"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

type GetListTeamItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	FoundedYear string `json:"founded_year"`
}

type GetTeamDetail struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Logo        string    `json:"logo"`
	FoundedYear string    `json:"founded_year"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDeleted   bool      `json:"is_deleted"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name"`
	Logo        *string `json:"logo"`
	FoundedYear *string `json:"founded_year"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
}
