package controllers

import (
	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
)

type IPlayersController interface {
	GetListPlayer(g *gin.Context) api.WebResponse[[]dto.GetListPlayerItem]
	UpdatePlayer(g *gin.Context) api.WebResponse[*dto.GetPlayerDetailResponse]
	GetPlayerDetail(g *gin.Context) api.WebResponse[*dto.GetPlayerDetailResponse]
	DeletePlayer(g *gin.Context) api.WebResponse[any]
}

type playerController struct {
	playerService services.IPlayersService
}

func NewPlayerController(
	playerService services.IPlayersService,
) IPlayersController {
	return &playerController{
		playerService: playerService,
	}
}

// GetListPlayer
//
//	@Summary		Get List Player
//	@Description	Get List Player
//	@Tags			Players
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.WebResponse[[]dto.GetListPlayerItem]
//	@Failure		500	{object}	api.WebResponse[any]
//	@Router			/api/v1/players [get]
func (c *playerController) GetListPlayer(g *gin.Context) api.WebResponse[[]dto.GetListPlayerItem] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetListPlayer].ctrl: Started").Msg()
	return c.playerService.GetListPlayer(ctx)
}

// UpdatePlayer
//
//	@Summary		Update Player Data
//	@Description	Update Player Data
//	@Tags			Players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string					true	"Player Id"
//	@Param			request		body		dto.UpdatePlayerRequest	true	"Update player request"
//	@Success		201			{object}	api.WebResponse[dto.GetPlayerDetailResponse]
//	@Failure		400			{object}	api.WebResponse[any]
//	@Failure		500			{object}	api.WebResponse[any]
//	@Router			/api/v1/players/{playerId} [put]
func (c *playerController) UpdatePlayer(g *gin.Context) api.WebResponse[*dto.GetPlayerDetailResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	var request dto.UpdatePlayerRequest
	if err := g.ShouldBindBodyWithJSON(&request); err != nil {
		l.WithContext(ctx).Error("error when read request").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.GetPlayerDetailResponse](
			ctx,
			constants.BadRequestDefault.GetMessage(),
			constants.BadRequestDefault.GetCode(),
			constants.BadRequestDefault.GetHttpCode(),
			nil,
		)
	}

	playerId := g.Param("playerId")
	return c.playerService.UpdatePlayer(ctx, playerId, &request)
}

// GetPlayerDetail
//
//	@Summary		Get Player Detail
//	@Description	Get Player Detail
//	@Tags			Players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string	true	"Player Id"
//	@Success		200			{object}	api.WebResponse[dto.GetPlayerDetailResponse]
//	@Failure		400			{object}	api.WebResponse[any]
//	@Failure		500			{object}	api.WebResponse[any]
//	@Router			/api/v1/players/{playerId} [get]
func (c *playerController) GetPlayerDetail(g *gin.Context) api.WebResponse[*dto.GetPlayerDetailResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetPlayerDetail].ctrl: Started").Msg()
	playerId := g.Param("playerId")
	return c.playerService.GetPlayerDetail(ctx, playerId)
}

// DeletePlayer
//
//	@Summary		Delete Player
//	@Description	Delete Player
//	@Tags			Players
//	@Accept			json
//	@Produce		json
//	@Param			playerId	path		string	true	"Player Id"
//	@Success		202			{object}	api.WebResponse[any]
//	@Failure		400			{object}	api.WebResponse[any]
//	@Failure		500			{object}	api.WebResponse[any]
//	@Router			/api/v1/players/{playerId} [delete]
func (c *playerController) DeletePlayer(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetPlayerDetail].ctrl: Started").Msg()
	playerId := g.Param("playerId")
	return c.playerService.DeletePlayer(ctx, playerId)
}
