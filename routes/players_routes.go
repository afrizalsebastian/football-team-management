package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupPlayerRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	playerGroup := router.Group("/players")
	{
		playerGroup.GET("/", server.GetListPlayer)
		playerGroup.PUT("/:playerId", server.UpdatePlayer)
	}
}
