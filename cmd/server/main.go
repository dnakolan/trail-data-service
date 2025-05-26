package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dnakolan/trail-data-service/internal/config"
	"github.com/dnakolan/trail-data-service/internal/handlers"
	"github.com/dnakolan/trail-data-service/internal/middleware"
	"github.com/dnakolan/trail-data-service/internal/services"
	"github.com/dnakolan/trail-data-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	router := gin.Default()
	gin.SetMode(cfg.Server.GinMode)

	trailsStorage := storage.NewTrailStorage()

	trailsService := services.NewTrailsService(trailsStorage)
	loginService := services.NewLoginService()

	healthHandler := handlers.NewHealthHandler()
	trailsHandler := handlers.NewTrailsHandler(trailsService)
	loginHandler := handlers.NewLoginHandler(loginService)

	router.POST("/login", loginHandler.LoginHandler)
	router.GET("/health", healthHandler.GetHealthHandler)

	router.POST("/trails", middleware.JwtAuthMiddleware(), trailsHandler.CreateTrailHandler)
	router.GET("/trails/:uid", middleware.JwtAuthMiddleware(), trailsHandler.GetTrailsHandler)
	router.GET("/trails", middleware.JwtAuthMiddleware(), trailsHandler.ListTrailsHandler)
	router.GET("/trails/nearby", middleware.JwtAuthMiddleware(), trailsHandler.ListTrailsHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("Received terminate, graceful shutdown", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server exited properly")

	os.Exit(0)
}
