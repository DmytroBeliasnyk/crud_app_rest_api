package main

import (
	"log"
	"os"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/datasource"
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

	cfg := datasource.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err := datasource.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("error occurred while connecting to db: %s", err.Error())
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("error occurred while close db: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
