package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupMatchRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	matchGroup := router.Group("/matches")
	{
		matchGroup.POST("/", server.CreateMatch)
		matchGroup.GET("/", server.GetListMatches)
		matchGroup.GET("/:matchId", server.MatchGoalList)
		matchGroup.DELETE("/:matchId", server.DeleteMatch)
		matchGroup.POST("/:matchId/goals", server.CreateMatchGoal)
		matchGroup.PUT("/:matchId/reschedule", server.RescheduleMatch)
		matchGroup.POST("/:matchId/full-time", server.MatchFullTime)
	}
}
