package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/handlers"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("error occurred while connecting to db: %s", err.Error())
	}

	repo := repositories.NewRepository(db)
	service := services.NewService(repo)
	handlers := handlers.NewHandler(service)

	server := new(core.Server)
	go func() {
		if err = server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
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

func initDB() (*sqlx.DB, error) {
	cfg := repositories.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	return repositories.NewPostgresDB(cfg)
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
