package app

import (
	"app/config"
	httpapi "app/internal/controller/http/v1"
	"app/pkg/httpserver"
	"app/pkg/postgres"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)


func Run() {
	// Config
	configPath, ok := os.LookupEnv("APP_CONFIG_PATH")
	if !ok || len(configPath) == 0{
		log.Fatal("app - os.LookupEnv: APP_CONFIG_PATH not specified")
	}

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("app - config.New: %w", err))
	}

	// Logger
	setLogrus(cfg.Log.Level)
	log.Info("Config read successfully...")

	// Postgres
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - postgres.New: %w", err))
	}
	defer pg.Close()

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	httpapi.ConfigureRouter(handler)

	// HttpServer
	log.Info("Starting HTTP server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Finish waiting
	log.Info("Configuring graceful shutdown...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <- signalChan:
		log.Infof("app - Run - signal: %s", s)
	case err := <- httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Graceful shutdown...")
	if err := httpServer.Shutdown(); err != nil {
		log.Error(fmt.Errorf("app - Run - httpSever.Shutdown: %w", err))
	}
}
