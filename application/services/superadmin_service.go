package services

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"golang.org/x/crypto/bcrypt"
)

type ISuperAdminService interface {
	GetHelloSuperAdmin(ctx context.Context) api.WebResponse[any]
	CreateAdminAccount(ctx context.Context, request *dto.CreateAdminAccount) api.WebResponse[any]
}

type superAdminService struct {
	adminAccountRepository repository.IAdminAccountRepository
}

func NewSuperAdminService(
	adminAccountRepository repository.IAdminAccountRepository,
) ISuperAdminService {
	return &superAdminService{
		adminAccountRepository: adminAccountRepository,
	}
}

func (h *superAdminService) GetHelloSuperAdmin(ctx context.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	l.WithContext(ctx).Info("Hello from service super admin").Msg()
	return api.SuccessResponse[any](ctx, constants.SuccessHello.GetMessage(), constants.SuccessHello.GetCode(), constants.SuccessHello.GetHttpCode(), nil)
}

func (h *superAdminService) CreateAdminAccount(ctx context.Context, request *dto.CreateAdminAccount) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateAdminAccount].service: Started").Msg()
	acct := &dao.AdminAccount{
		Username: request.Username,
		Password: request.Password,
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(acct.Password), bcrypt.DefaultCost)
	if err != nil {
		l.WithContext(ctx).Error("failed to hash password").Attr("err", err).Msg()
		return api.ErrorResponse[any](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}
	acct.Password = string(hashedPass)

	if err := h.adminAccountRepository.CreateAccount(ctx, acct); err != nil {
		l.WithContext(ctx).Error("failed to create account").Attr("err", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.DuplicateRow) {
			errMsg = constants.BadRequestDefault
		}

		return api.ErrorResponse[any](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	l.WithContext(ctx).Debug("[CreateAdminAccount].service: Completed").Msg()
	return api.ErrorResponse[any](
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		nil,
	)
}
