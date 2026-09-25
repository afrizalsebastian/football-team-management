package controllers

import (
	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
)

type ISuperAdminController interface {
	GetHelloSuperAdmin(g *gin.Context) api.WebResponse[any]
	CreateAccount(g *gin.Context) api.WebResponse[any]
}

type superAdminController struct {
	superAdminService services.ISuperAdminService
}

func NewSuperAdminController(
	superAdminService services.ISuperAdminService,
) ISuperAdminController {
	return &superAdminController{
		superAdminService: superAdminService,
	}
}

// GetHelloSuperAdmin
//
//	@Summary		Example API -- Hello from super admin
//	@Description	Example API
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//
//	@Param			x-internal-token	header		string	true	"token"
//
//	@Success		200					{object}	api.WebResponse[any]
//	@Success		401					{object}	api.WebResponse[any]
//	@Router			/api/v1/internal/admin [get]
func (c *superAdminController) GetHelloSuperAdmin(g *gin.Context) api.WebResponse[any] {
	ctx := g.Request.Context()
	return c.superAdminService.GetHelloSuperAdmin(ctx)
}

// CreateAccount
//
//	@Summary		CreateAccount
//	@Description	CreateAccount
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//
//	@Param			x-internal-token	header		string					true	"token"
//
//	@Param			{object}			body		dto.CreateAdminAccount	true	"create account request"
//
//	@Success		200					{object}	api.WebResponse[any]
//	@Success		400					{object}	api.WebResponse[any]
//	@Success		401					{object}	api.WebResponse[any]
//	@Success		500					{object}	api.WebResponse[any]
//	@Router			/api/v1/internal/admin [post]
func (c *superAdminController) CreateAccount(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[CreateAccount].ctrl: Started").Msg()

	var request dto.CreateAdminAccount
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[any](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	if valErr := helper.ValidateParams(ctx, &request); valErr != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", valErr).Msg()
		return api.ErrorResponse[any](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			valErr,
		)
	}

	return c.superAdminService.CreateAdminAccount(ctx, &request)
}
