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
	HelloController  controllers.IHelloController
	TeamsController  controllers.ITeamsController
	PlayerController controllers.IPlayersController
	MatchController  controllers.IMatchController
	GoalController   controllers.IGoalController
}

func NewServer(app *bootstrap.FootballManagementApp) *HttpServer {
	di := initDI(app)
	return &HttpServer{
		HelloController:  di.HelloController,
		TeamsController:  di.TeamsController,
		PlayerController: di.PlayerController,
		MatchController:  di.MatchController,
		GoalController:   di.GoalController,
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

func (s *HttpServer) GetListTeam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.GetListTeam(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) SoftDeleteTeam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.SoftDeleteTeam(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) CreateTeamPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.CreateTeamPlayer(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) GetListTeamPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.GetListTeamPlayer(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) GetTeamDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.GetTeamDetail(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) UpdateTeam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.TeamsController.UpdateTeam(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) GetListPlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.PlayerController.GetListPlayer(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) UpdatePlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.PlayerController.UpdatePlayer(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) GetPlayerDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.PlayerController.GetPlayerDetail(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) DeletePlayer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.PlayerController.DeletePlayer(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) CreateMatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.CreateMatch(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) GetListMatches(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.GetListMatches(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) RescheduleMatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.RescheduleMatch(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) DeleteMatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.DeleteMatch(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) CreateMatchGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.CreateMatchGoal(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) MatchGoalList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.MatchController.MatchGoalList(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}

func (s *HttpServer) DeleteGoal(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)
	resp := s.GoalController.DeleteGoal(c)
	api.WriteJSONResponse(c, resp.HttpCode, resp)
}
