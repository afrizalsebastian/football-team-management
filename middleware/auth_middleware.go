package middleware

import (
	"strings"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/bootstrap"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(app *bootstrap.FootballManagementApp) gin.HandlerFunc {
	l := logger.LoggerNew()
	return func(g *gin.Context) {
		headerValue := g.GetHeader("Authorization")
		ctx := g.Request.Context()
		if headerValue == "" || !strings.HasPrefix(headerValue, "Bearer ") {
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

		token := strings.TrimPrefix(headerValue, "Bearer ")
		claims, err := helper.VerifyToken(token, app.Env.JwtSecret)
		if err != nil {
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

		ctx = helper.SetContextValueClaims(ctx, claims)
		g.Request = g.Request.WithContext(ctx)
		g.Next()
	}
}
