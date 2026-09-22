package dto

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
