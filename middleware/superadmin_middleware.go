package middleware

import (
	"slices"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
)

func SuperadminMiddleware(app *bootstrap.FootballManagementApp) gin.HandlerFunc {
	l := logger.LoggerNew()
	return func(g *gin.Context) {
		token := g.GetHeader("x-internal-token")
		ctx := g.Request.Context()
		if token == "" ||
			len(app.Env.SuperadminToken) == 0 ||
			!slices.Contains(app.Env.SuperadminToken, token) {
			l.WithContext(ctx).Warn("invalid token").Attr("token", token).Msg()
			resp := api.ErrorResponse[any](
				ctx,
				constants.UnauthorizedDefault.GetMessage(),
				constants.UnauthorizedDefault.GetCode(),
				constants.UnauthorizedDefault.GetHttpCode(),
				nil,
			)
			api.WriteJSONResponse(g, resp.HttpCode, resp)
			g.Abort()
			return
		}
		g.Next()
	}
}
