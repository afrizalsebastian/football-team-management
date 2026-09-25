package controllers

import (
	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
)

type IGoalController interface {
	DeleteGoal(g *gin.Context) api.WebResponse[any]
}

type goalController struct {
	goalService services.IGoalService
}

func NewGoalController(
	goalService services.IGoalService,
) IGoalController {
	return &goalController{
		goalService: goalService,
	}
}

// DeleteGoal
//
//	@Summary		Delete goal (soft-delete)
//	@Description	Delete goal (soft-delete)
//	@Tags			Goals
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			goalId	path		string	true	"Goal ID"
//
//	@Success		202		{object}	api.WebResponse[any]
//	@Failure		400		{object}	api.WebResponse[any]
//	@Failure		401		{object}	api.WebResponse[any]
//	@Failure		500		{object}	api.WebResponse[any]
//	@Router			/api/v1/goals/{goalId} [delete]
func (c *goalController) DeleteGoal(g *gin.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	ctx := g.Request.Context()

	l.WithContext(ctx).Debug("[DeleteGoal].ctrl: Started").Msg()
	goalId := g.Param("goalId")
	return c.goalService.DeleteGoals(ctx, goalId)
}
