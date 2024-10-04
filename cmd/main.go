package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/handlers"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
)

const (
	CONFIG_FOLDER = "configs"
	CONFIG_FILE   = "main"
)

//	@title		Documentation for api
//	@version	1.0

//	@host		localhost:8000
//	@BasePath	/

// @accept		json
// @produce	json
func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	db, err := repositories.NewPostgresDB(*cfg)
	if err != nil {
		log.Fatalf("error occurred while connecting to db: %s", err.Error())
	}

	repo := repositories.NewRepository(db)
	service := services.NewService(repo)
	handlers := handlers.NewHandler(service)

	server := new(core.Server)
	go func() {
		if err = server.Run(cfg.ServerPort, handlers.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err = server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occurred on server shutting down: %s", err.Error())
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("error occurred on db connection close: %s", err.Error())
	}
}

func initConfig() (*config.Config, error) {
	cfg, err := config.InitConfig(CONFIG_FOLDER, CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
