package controllers

import (
	"strings"

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
	GetListMatches(g *gin.Context) api.WebResponse[[]dto.GetMatchResponse]
	RescheduleMatch(g *gin.Context) api.WebResponse[any]
	DeleteMatch(g *gin.Context) api.WebResponse[any]
	CreateMatchGoal(g *gin.Context) api.WebResponse[any]
	MatchGoalList(g *gin.Context) api.WebResponse[*dto.MatchGoalsResponse]
	MatchFullTime(g *gin.Context) api.WebResponse[any]
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
//
//	@Security		BearerAuth
//
//	@Param			request	body		dto.CreateMatchRequest	true	"Create match body request"
//	@Success		201		{object}	api.WebResponse[dto.CreateMatchResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
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

	if helper.MatchDateTimeInPast(request.Date, request.Time) {
		return api.ErrorResponse[*dto.CreateMatchResponse](
			ctx,
			constants.MatchDateInPast.GetMessage(),
			constants.MatchDateInPast.GetCode(),
			constants.MatchDateInPast.GetHttpCode(),
			nil,
		)
	}

	return c.matchService.CreateMatch(ctx, &request)
}

// GetListMatches
//
//	@Summary		Get List Matches
//	@Description	Get List Matches
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Param			teamId	query		string	false	"Team ID"
//
//	@Param			type	query		string	false	"type match: one of 'result' or 'scheduled'"
//
//	@Success		200		{object}	api.WebResponse[[]dto.GetMatchResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches [get]
func (c *matchController) GetListMatches(g *gin.Context) api.WebResponse[[]dto.GetMatchResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[GetListMatches].ctrl: Started").Msg()

	var teamId *string
	if val := g.Query("teamId"); val != "" {
		teamId = helper.StringPtr(val)
	}

	var matchType *string
	if val := g.Query("type"); val != "" && (strings.EqualFold(val, "result") || strings.EqualFold(val, "scheduled")) {
		matchType = helper.StringPtr(strings.ToLower(val))
	}

	return c.matchService.GetListMatches(ctx, teamId, matchType)
}

// RescheduleMatch
//
//	@Summary		Reschedule Match
//	@Description	Reschedule match
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			matchId	path		string						true	"Match ID"
//
//	@Param			request	body		dto.RescheduleMatchRequest	true	"Reschedule match request"
//	@Success		202		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches/{matchId}/reschedule [put]
func (c *matchController) RescheduleMatch(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[RescheduleMatch].ctrl: Started").Msg()

	var request dto.RescheduleMatchRequest
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

	if helper.MatchDateTimeInPast(request.Date, request.Time) {
		return api.ErrorResponse[any](
			ctx,
			constants.MatchDateInPast.GetMessage(),
			constants.MatchDateInPast.GetCode(),
			constants.MatchDateInPast.GetHttpCode(),
			nil,
		)
	}

	matchId := g.Param("matchId")
	return c.matchService.RescheduleMatch(ctx, matchId, &request)
}

// DeleteMatch
//
//	@Summary		Delete Match
//	@Description	Delete Match
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			matchId	path		string	true	"Match ID"
//
//	@Success		200		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
//	@Failure		404		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches/{matchId} [delete]
func (c *matchController) DeleteMatch(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[DeleteMatch].ctrl: Started").Msg()
	matchId := g.Param("matchId")
	return c.matchService.DeleteMatch(ctx, matchId)
}

// CreateMatchGoal
//
//	@Summary		Create Match goal
//	@Description	Create a new match goal
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			matchId	path		string			true	"Match ID"
//
//	@Param			request	body		dto.MatchGoal	true	"Create match body request"
//	@Success		201		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
//	@Failure		404		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches/{matchId}/goals [post]
func (c *matchController) CreateMatchGoal(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[CreateMatchGoal].ctrl: Started").Msg()

	var request dto.MatchGoal
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

	matchId := g.Param("matchId")
	return c.matchService.CreateMatchGoal(ctx, matchId, &request)
}

// MatchGoalList
//
//	@Summary		Get Match Goal list
//	@Description	Get Match Goal list
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Param			matchId	path		string	true	"Team ID"
//
//	@Success		200		{object}	api.WebResponse[dto.MatchGoalsResponse]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches/{matchId} [get]
func (c *matchController) MatchGoalList(g *gin.Context) api.WebResponse[*dto.MatchGoalsResponse] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[MatchGoalList].ctrl: Started").Msg()

	matchId := g.Param("matchId")
	return c.matchService.MatchGoalList(ctx, matchId)
}

// MatchFullTime
//
//	@Summary		Submit matches result. Full-Time
//	@Description	Submit matches result. Full-Time
//	@Tags			Matches
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			matchId	path		string	true	"Match ID"
//
//	@Success		200		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
//	@Failure		404		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/matches/{matchId}/full-time [post]
func (c *matchController) MatchFullTime(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[MatchFullTime].ctrl: Started").Msg()
	matchId := g.Param("matchId")
	return c.matchService.MatchFullTime(ctx, matchId)
}
