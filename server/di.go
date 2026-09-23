package server

import (
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
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
	teamRepository := repository.NewTeamsRepository(app.DBPool)
	playerRepository := repository.NewPlayerRepository(app.DBPool)
	service := services.NewTeamsService(teamRepository, playerRepository)
	ctrl := controllers.NewTeamsController(service)

	return ctrl
}
