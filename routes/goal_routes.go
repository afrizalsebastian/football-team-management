package routes

import (
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-gonic/gin"
)

func SetupGoalRoutes(router *gin.RouterGroup, server *server.HttpServer) {
	goalGroup := router.Group("/goals")
	{
		goalGroup.DELETE("/:goalId", server.DeleteGoal)
	}
}
