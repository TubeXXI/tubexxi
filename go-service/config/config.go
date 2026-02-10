package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Scraper  ScraperConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ScraperConfig struct {
	Host string
	Port string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var config Config

	config.App.Name = viper.GetString("APP_NAME")
	config.App.Port = viper.GetString("APP_PORT")
	config.App.Env = viper.GetString("APP_ENV")

	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetString("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.Name = viper.GetString("DB_NAME")

	config.Redis.Host = viper.GetString("REDIS_HOST")
	config.Redis.Port = viper.GetString("REDIS_PORT")
	config.Redis.Password = viper.GetString("REDIS_PASSWORD")
	config.Redis.DB = viper.GetInt("REDIS_DB")

	config.Scraper.Host = viper.GetString("SCRAPER_HOST")
	config.Scraper.Port = viper.GetString("SCRAPER_PORT")

	return &config
}
