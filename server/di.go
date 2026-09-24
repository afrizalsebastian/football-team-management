package server

import (
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
)

type ServerDependencies struct {
	HelloController  controllers.IHelloController
	TeamsController  controllers.ITeamsController
	PlayerController controllers.IPlayersController
	MatchController  controllers.IMatchController
	GoalController   controllers.IGoalController
}

func initDI(app *bootstrap.FootballManagementApp) *ServerDependencies {
	return &ServerDependencies{
		HelloController:  setupHelloControllers(),
		TeamsController:  setupTeamsController(app),
		PlayerController: setupPlayerController(app),
		MatchController:  setupMatchController(app),
		GoalController:   setupGoalController(app),
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

func setupPlayerController(app *bootstrap.FootballManagementApp) controllers.IPlayersController {
	teamRepository := repository.NewTeamsRepository(app.DBPool)
	playerRepository := repository.NewPlayerRepository(app.DBPool)
	service := services.NewPlayerService(playerRepository, teamRepository)
	ctrl := controllers.NewPlayerController(service)

	return ctrl
}

func setupMatchController(app *bootstrap.FootballManagementApp) controllers.IMatchController {
	matchRepository := repository.NewMatchesRepository(app.DBPool)
	goalRepository := repository.NewGoalsRepository(app.DBPool)
	service := services.NewMatchService(matchRepository, goalRepository)
	ctrl := controllers.NewMatchController(service)

	return ctrl
}

func setupGoalController(app *bootstrap.FootballManagementApp) controllers.IGoalController {
	goalRepository := repository.NewGoalsRepository(app.DBPool)
	service := services.NewGoalService(goalRepository)
	ctrl := controllers.NewGoalController(service)

	return ctrl
}
