package services

import (
	"context"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/module/logger"
)

type IHelloService interface {
	GetHello(ctx context.Context) api.WebResponse[any]
}

type helloService struct{}

func NewHelloService() IHelloService {
	return &helloService{}
}

func (h *helloService) GetHello(ctx context.Context) api.WebResponse[any] {
	l := logger.LoggerNew()
	l.WithContext(ctx).Info("Hello from service").Msg()
	return api.SuccessResponse[any](ctx, constants.SuccessHello.GetMessage(), constants.SuccessHello.GetCode(), constants.SuccessHello.GetHttpCode(), nil)
}
