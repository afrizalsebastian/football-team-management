package dto

import "github.com/golang-jwt/jwt/v5"

type CreateAdminAccount struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminClaims struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`

	jwt.RegisteredClaims
}
type LoginRequest struct {
	CreateAdminAccount
}

type LoginResponse struct {
	Token string `json:"token"`
}
