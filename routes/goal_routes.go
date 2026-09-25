package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupGoalRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	goalGroup := router.Group("/goals")
	{
		goalGroup.DELETE("/:goalId", middleware.AuthMiddleware(app), server.DeleteGoal)
	}
}
