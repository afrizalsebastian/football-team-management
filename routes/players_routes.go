package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupPlayerRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	playerGroup := router.Group("/players")
	{
		playerGroup.GET("/", server.GetListPlayer)
		playerGroup.GET("/:playerId", server.GetPlayerDetail)
		playerGroup.PUT("/:playerId", middleware.AuthMiddleware(app), server.UpdatePlayer)
		playerGroup.DELETE("/:playerId", middleware.AuthMiddleware(app), server.DeletePlayer)
	}
}
