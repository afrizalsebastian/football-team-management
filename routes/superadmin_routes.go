package routes

import (
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupSuperAdminRoutes(router *gin.RouterGroup, app *bootstrap.FootballManagementApp, server *server.HttpServer) {
	superadminGroup := router.Group("/internal/admin")
	superadminGroup.Use(middleware.SuperadminMiddleware(app))
	{
		superadminGroup.GET("/", server.GetHelloSuperAdmin)
		superadminGroup.POST("/", server.CreateAccount)
	}
}
