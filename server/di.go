package server

import (
	"github.com/afrizalsebastian/football-team-management/application/controllers"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
)

type ServerDependencies struct {
	HelloController controllers.IHelloController
}

func initDI(app *bootstrap.FootballManagementApp) *ServerDependencies {
	return &ServerDependencies{
		HelloController: setupHelloControllers(),
	}
}

func setupHelloControllers() controllers.IHelloController {
	svc := services.NewHelloService()
	ctrl := controllers.NewHelloController(svc)

	return ctrl
}
