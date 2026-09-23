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

type IMatchController interface {
	CreateMatch(g *gin.Context) api.WebResponse[*dto.CreateMatchResponse]
}

type matchController struct {
	matchService services.IMatchesService
}

func NewMatchController(
	matchService services.IMatchesService,
) IMatchController {
	return &matchController{
		matchService: matchService,
	}
}

// CreateMatch
//
//	@Summary		Create Match
//	@Description	Create a new match
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateMatchRequest	true	"Create match body request"
//	@Success		201		{object}	api.WebResponse[dto.CreateMatchResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches [post]
func (c *matchController) CreateMatch(g *gin.Context) api.WebResponse[*dto.CreateMatchResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[CreateMatch].ctrl: Started").Msg()

	var request dto.CreateMatchRequest
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.CreateMatchResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	if valErr := helper.ValidateParams(ctx, &request); valErr != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", valErr).Msg()
		return api.ErrorResponse[*dto.CreateMatchResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			valErr,
		)
	}

	return c.matchService.CreateMatch(ctx, &request)
}
