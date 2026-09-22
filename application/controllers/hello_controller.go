package controllers

import (
	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/services"
	"github.com/gin-gonic/gin"
)

type IHelloController interface {
	GetHello(g *gin.Context) api.WebResponse[any]
}

type helloController struct {
	helloService services.IHelloService
}

func NewHelloController(
	helloService services.IHelloService,
) IHelloController {
	return &helloController{
		helloService: helloService,
	}
}

// GetHello
//
//	@Summary		Example API -- Hello
//	@Description	Example API
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.WebResponse[any]
//	@Router			/api/v1/hello [get]
func (h *helloController) GetHello(g *gin.Context) api.WebResponse[any] {
	ctx := g.Request.Context()
	return h.helloService.GetHello(ctx)
}
