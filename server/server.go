package server

import (
	"context"
	"time"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	HelloController controllers.IHelloController
	TeamsController controllers.ITeamsController
}

func NewServer(app *bootstrap.FootballManagementApp) *HttpServer {
	di := initDI(app)
	return &HttpServer{
		HelloController: di.HelloController,
		TeamsController: di.TeamsController,
	}
}

func (s *HttpServer) GetHello(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.HelloController.GetHello(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) PostCreateTeam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.CreateTeam(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}
