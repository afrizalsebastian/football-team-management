package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupTeamsRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	teamGroup := router.Group("/teams")
	{
		teamGroup.POST("/", server.PostCreateTeam)
		teamGroup.GET("/", server.GetListTeam)
		teamGroup.DELETE("/:teamId", server.SoftDeleteTeam)
		teamGroup.POST("/:teamId/players", server.CreateTeamPlayer)
	}
}
