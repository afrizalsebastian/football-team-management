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

type IAdminController interface {
	Login(g *gin.Context) api.WebResponse[*dto.LoginResponse]
}

type adminController struct {
	adminService services.IAdminService
}

func NewAdminController(
	adminService services.IAdminService,
) IAdminController {
	return &adminController{
		adminService: adminService,
	}
}

// Login
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//
//	@Param			{object}	body		dto.LoginRequest	true	"login request"
//
//	@Success		200			{object}	api.WebResponse[any]
//	@Success		400			{object}	api.WebResponse[any]
//	@Success		500			{object}	api.WebResponse[any]
//	@Router			/api/v1/admin [post]
func (c *adminController) Login(g *gin.Context) api.WebResponse[*dto.LoginResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[Login].ctrl: Started").Msg()

	var request dto.LoginRequest
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.LoginResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	if valErr := helper.ValidateParams(ctx, &request); valErr != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", valErr).Msg()
		return api.ErrorResponse[*dto.LoginResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			valErr,
		)
	}

	return c.adminService.Login(ctx, &request)
}
