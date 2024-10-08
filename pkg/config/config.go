package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"server_port"`

	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		Password DBPassword
	} `mapstructure:"db"`

	Auth struct {
		Salt      string
		Signature string
		JWT       time.Duration `mapstructure:"jwt"`
		Refresh   time.Duration `mapstructure:"refresh"`
	} `mapstructure:"tokens_ttl"`
}

type DBPassword struct {
	Password string
}

func InitConfig(folder, file string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := parseConfig(cfg); err != nil {
		return nil, err
	}

	if err := parseEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseConfig(cfg *Config) error {
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("db", cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("tokens_ttl", cfg); err != nil {
		return err
	}

	return nil
}

func parseEnv(cfg *Config) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := envconfig.Process("db", &cfg.DB.Password); err != nil {
		return err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return err
	}

	return nil
}
