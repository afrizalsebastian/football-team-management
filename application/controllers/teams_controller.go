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

type ITeamsController interface {
	CreateTeam(g *gin.Context) api.WebResponse[*dto.CreateTeamResponse]
	GetListTeam(g *gin.Context) api.WebResponse[[]dto.GetListTeamItem]
	SoftDeleteTeam(g *gin.Context) api.WebResponse[any]
	CreateTeamPlayer(g *gin.Context) api.WebResponse[*dto.CreatePlayerTeamResponse]
}

type teamsController struct {
	teamsService services.ITeamsService
}

func NewTeamsController(teamsService services.ITeamsService) ITeamsController {
	return &teamsController{
		teamsService: teamsService,
	}
}

// CreateTeam
//
//	@Summary		Create Team
//	@Description	Create a new team
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateTeamRequest	true	"Create team bory request"
//	@Success		201		{object}	api.WebResponse[dto.CreateTeamResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/teams [post]
func (c *teamsController) CreateTeam(g *gin.Context) api.WebResponse[*dto.CreateTeamResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[CreateTeam].ctrl: Started").Msg()

	var request dto.CreateTeamRequest
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.CreateTeamResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	if valErr := helper.ValidateParams(ctx, &request); valErr != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", valErr).Msg()
		return api.ErrorResponse[*dto.CreateTeamResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			valErr,
		)
	}

	return c.teamsService.CreateTeam(ctx, &request)
}

// GetListTeam
//
//	@Summary		Get List team
//	@Description	Get List Team
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.WebResponse[[]dto.GetListTeamItem]
//	@Failure		500	{object}	api.WebResponse[any]
//	@Router			/api/v1/teams [get]
func (c *teamsController) GetListTeam(g *gin.Context) api.WebResponse[[]dto.GetListTeamItem] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetListTeam].ctrl: Started").Msg()
	return c.teamsService.GetListTeam(ctx)
}

// SoftDeleteTeam
//
//	@Summary		Delete team (soft-delete)
//	@Description	Delete team (soft-delete)
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//
//	@Param			teamId	path		string	true	"Team ID"
//
//	@Success		202		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/teams/{teamId} [delete]
func (c *teamsController) SoftDeleteTeam(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetListTeam].ctrl: Started").Msg()
	teamId := g.Param("teamId")
	return c.teamsService.SoftDeleteTeam(ctx, teamId)
}

// CreateTeamPlayer
//
//	@Summary		Create Team Player
//	@Description	Create a new team Player
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//
//	@Param			teamId	path		string						true	"Team ID"
//
//	@Param			request	body		dto.CreatePlayerTeamRequest	true	"Create team bory request"
//	@Success		201		{object}	api.WebResponse[dto.CreatePlayerTeamResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/teams/{teamId}/players [post]
func (c *teamsController) CreateTeamPlayer(g *gin.Context) api.WebResponse[*dto.CreatePlayerTeamResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[CreateTeamPlayer].ctrl: Started").Msg()

	var request dto.CreatePlayerTeamRequest
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	if valErr := helper.ValidateParams(ctx, &request); valErr != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", valErr).Msg()
		return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			valErr,
		)
	}

	teamId := g.Param("teamId")
	return c.teamsService.CreatePlayerTeam(ctx, teamId, &request)
}
