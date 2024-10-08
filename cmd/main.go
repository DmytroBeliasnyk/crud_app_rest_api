package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/handlers"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/implserv"
	"github.com/DmytroBeliasnyk/in_memory_cache/memory"
	"github.com/sirupsen/logrus"
)

const (
	CONFIG_FOLDER = "configs"
	CONFIG_FILE   = "main"
)

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.DateTime,
	})
}

//	@title		Documentation for api
//	@version	1.0

//	@host		localhost:8000
//	@BasePath	/

// @accept		json
// @produce	json
func main() {
	cfg, err := initConfig()
	if err != nil {
		logrus.WithField("error", err).Fatal("error initializing config")
	}

	db, err := repositories.NewPostgresDB(cfg)
	if err != nil {
		logrus.WithField("error", err).Fatal("error occurred while connecting to db")
	}

	repo := repositories.NewRepository(db)
	refresh := repositories.NewRefreshTokenRepository(db)
	auth := implserv.NewAuthService(refresh, cfg)
	service := services.NewService(repo, auth)
	handlers := handlers.NewHandler(service, auth, cfg, memory.GetCache())

	server := new(core.Server)
	go func() {
		if err = server.Run(cfg.ServerPort, handlers.InitRoutes()); err != nil {
			logrus.WithField("error", err).Fatal("error occurred while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err = server.Shutdown(context.Background()); err != nil {
		logrus.WithField("error", err).Fatal("error occurred on server shutting down")
	}

	err = db.Close()
	if err != nil {
		logrus.WithField("error", err).Fatal("error occurred on db connection close")
	}
}

func initConfig() (*config.Config, error) {
	cfg, err := config.InitConfig(CONFIG_FOLDER, CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
