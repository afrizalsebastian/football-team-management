package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	superadminGroup := router.Group("/admin")
	{
		superadminGroup.POST("/", server.Login)
	}
}
