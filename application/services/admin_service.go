package services

import (
	"context"
	"errors"
	"time"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAdminService interface {
	Login(ctx context.Context, request *dto.LoginRequest) api.WebResponse[*dto.LoginResponse]
}

type adminService struct {
	adminAccountRepository repository.IAdminAccountRepository
	jwtSecret              string
	jwtTTL                 time.Duration
}

func NewAdminService(
	adminAccountRepository repository.IAdminAccountRepository,
	jwtSecret string,
	jwtTTL time.Duration,
) IAdminService {
	return &adminService{
		adminAccountRepository: adminAccountRepository,
		jwtSecret:              jwtSecret,
		jwtTTL:                 jwtTTL,
	}
}

func (s *adminService) Login(ctx context.Context, request *dto.LoginRequest) api.WebResponse[*dto.LoginResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[Login].service: Started").Msg()

	admin, err := s.adminAccountRepository.GetAdminAccount(ctx, request.Username)
	if err != nil {
		l.WithContext(ctx).Error("error when get account").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.InvalidUsernameOrPassword
		}

		return api.ErrorResponse[*dto.LoginResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password)); err != nil {
		l.WithContext(ctx).Error("invalid password").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.LoginResponse](
			ctx,
			constants.InvalidUsernameOrPassword.GetMessage(),
			constants.InvalidUsernameOrPassword.GetCode(),
			constants.InvalidUsernameOrPassword.GetHttpCode(),
			nil,
		)
	}

	now := time.Now()
	claims := &dto.AdminClaims{
		Username: admin.Username,
		UserId:   admin.Id,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "football-management-app",
			Subject:   admin.Id,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
		},
	}

	token, err := helper.GenerateJWT(s.jwtSecret, claims)
	if err != nil {
		l.WithContext(ctx).Error("failed to generate token").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.LoginResponse](
			ctx,
			constants.InvalidUsernameOrPassword.GetMessage(),
			constants.InvalidUsernameOrPassword.GetCode(),
			constants.InvalidUsernameOrPassword.GetHttpCode(),
			nil,
		)
	}

	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		&dto.LoginResponse{
			Token: token,
		},
	)
}
