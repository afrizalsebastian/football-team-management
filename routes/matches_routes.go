package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupMatchRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	matchGroup := router.Group("/matches")
	{
		matchGroup.POST("/", middleware.AuthMiddleware(app), server.CreateMatch)
		matchGroup.GET("/", server.GetListMatches)
		matchGroup.GET("/:matchId", server.MatchGoalList)
		matchGroup.DELETE("/:matchId", middleware.AuthMiddleware(app), server.DeleteMatch)
		matchGroup.POST("/:matchId/goals", middleware.AuthMiddleware(app), server.CreateMatchGoal)
		matchGroup.PUT("/:matchId/reschedule", middleware.AuthMiddleware(app), server.RescheduleMatch)
		matchGroup.POST("/:matchId/full-time", middleware.AuthMiddleware(app), server.MatchFullTime)
	}
}
