package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupTeamsRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	helloGroup := router.Group("/teams")
	{
		helloGroup.POST("/", server.PostCreateTeam)
	}
}
