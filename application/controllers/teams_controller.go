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

	l.WithContext(ctx).Debug("CreateTeam").Msg()

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
