package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupTeamsRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	teamGroup := router.Group("/teams")
	{
		teamGroup.POST("/", middleware.AuthMiddleware(app), server.PostCreateTeam)
		teamGroup.GET("/", server.GetListTeam)
		teamGroup.GET("/:teamId", server.GetTeamDetail)
		teamGroup.PUT("/:teamId", middleware.AuthMiddleware(app), server.UpdateTeam)
		teamGroup.DELETE("/:teamId", middleware.AuthMiddleware(app), server.SoftDeleteTeam)
		teamGroup.POST("/:teamId/players", middleware.AuthMiddleware(app), server.CreateTeamPlayer)
		teamGroup.GET("/:teamId/players", server.GetListTeamPlayer)
	}
}
