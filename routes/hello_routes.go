package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupHelloRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	helloGroup := router.Group("/hello")
	{
		helloGroup.GET("/", server.GetHello)
	}
}
