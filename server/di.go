package server

import (
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
)

type ServerDependencies struct {
	HelloController      controllers.IHelloController
	TeamsController      controllers.ITeamsController
	PlayerController     controllers.IPlayersController
	MatchController      controllers.IMatchController
	GoalController       controllers.IGoalController
	SuperAdminController controllers.ISuperAdminController
	AdminController      controllers.IAdminController
}

func initDI(app *bootstrap.FootballManagementApp) *ServerDependencies {
	return &ServerDependencies{
		HelloController:      setupHelloControllers(),
		TeamsController:      setupTeamsController(app),
		PlayerController:     setupPlayerController(app),
		MatchController:      setupMatchController(app),
		GoalController:       setupGoalController(app),
		SuperAdminController: setupSuperAdminControllers(app),
		AdminController:      setupAdminControllers(app),
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

func setupSuperAdminControllers(app *bootstrap.FootballManagementApp) controllers.ISuperAdminController {
	adminAccountRepository := repository.NewAdminAccountRepository(app.DBPool)
	svc := services.NewSuperAdminService(adminAccountRepository)
	ctrl := controllers.NewSuperAdminController(svc)

	return ctrl
}

func setupAdminControllers(app *bootstrap.FootballManagementApp) controllers.IAdminController {
	adminAccountRepository := repository.NewAdminAccountRepository(app.DBPool)
	svc := services.NewAdminService(adminAccountRepository, app.Env.JwtSecret, app.Env.JwtTTL)
	ctrl := controllers.NewAdminController(svc)

	return ctrl
}
