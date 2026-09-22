package middleware

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RecoveryPanicMiddleware() gin.HandlerFunc {
	l := logger.LoggerNew()
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				l.Info("Recoverd from panic").Attr("rec", rec).Msg()
				ctx.JSON(
					http.StatusInternalServerError,
					"Internal Server Error",
				)
			}
		}()
		ctx.Next()
	}
}

func RequestTracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		msgId := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), logger.ContextKeyMessageId, msgId)
		c.Header("X-Message-ID", msgId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
