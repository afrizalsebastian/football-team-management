package server

import (
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	teams_repository "github.com/afrizalsebastian/football-team-management/domain/repository/teams-repository"
)

type ServerDependencies struct {
	HelloController controllers.IHelloController
	TeamsController controllers.ITeamsController
}

func initDI(app *bootstrap.FootballManagementApp) *ServerDependencies {
	return &ServerDependencies{
		HelloController: setupHelloControllers(),
		TeamsController: setupTeamsController(app),
	}
}

func setupHelloControllers() controllers.IHelloController {
	svc := services.NewHelloService()
	ctrl := controllers.NewHelloController(svc)

	return ctrl
}

func setupTeamsController(app *bootstrap.FootballManagementApp) controllers.ITeamsController {
	repository := teams_repository.NewTeamsRepository(app.DBPool)
	service := services.NewTeamsService(repository)
	ctrl := controllers.NewTeamsController(service)

	return ctrl
}
