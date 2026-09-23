package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/football-team-management/bootstrap"
	_ "github.com/afrizalsebastian/football-team-management/docs"
	"github.com/afrizalsebastian/football-team-management/middleware"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/afrizalsebastian/football-team-management/routes"
	"github.com/afrizalsebastian/football-team-management/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			API Football Team Management
//	@version		1.0
//	@description	API Spec for footbal team management
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	afrizalsebastian
//	@contact.email	sebastiangurning@gmail.com

// @host	localhost:8080
func main() {
	l := logger.LoggerNew()

	app := bootstrap.NewApplication()
	r := gin.Default()
	r.Use(middleware.RecoveryPanicMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "404 -- route not found",
		})
	})

	httpServer := server.NewServer(app)
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.RequestTracingMiddleware())
	{
		routes.SetupHelloRoutes(apiV1, httpServer)
		routes.SetupTeamsRoutes(apiV1, httpServer)
		routes.SetupPlayerRoutes(apiV1, httpServer)
		routes.SetupMatchRoutes(apiV1, httpServer)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.Env.ServerPort),
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start server
	go func() {
		l.Infof("🚀 Server is running on port %v", app.Env.ServerPort).Msg()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("server failed to start").Msg()
		}
	}()

	<-signalChan
	l.Info("Shutting down server...").Msg()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Error("Server forced to shutdown").Msg()
	}

	l.Info("✅ Server exited gracefully").Msg()
}
