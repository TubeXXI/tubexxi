package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Centrifugo CentrifugoConfig
	MinIO      MinIOConfig
	Bycrypt    BycryptConfig
	Telegram   TelegramConfig
	Email      EmailConfig
	Scraper    ScraperConfig
}
type AppConfig struct {
	AppName    string
	AppEnv     string
	IsDebug    bool
	Port       string
	URL        string
	ClientUrl  string
	AdminEmail string
}

type DatabaseConfig struct {
	DbHost            string
	DbPort            string
	DbUser            string
	DbPassword        string
	DbName            string
	DbSSLMode         string
	DbMaxOpenConn     int
	DbMaxIdleConn     int
	DbConnMaxLifetime time.Duration
}

type RedisConfig struct {
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	RedisAsynqDB     int
	RedisInstance    string
	RedisPoolSize    int
	RedisConcurrency int
}

type JWTConfig struct {
	JwtSecret            string
	JwtExpiration        time.Duration
	JwtRefreshExpiration time.Duration
}

type CentrifugoConfig struct {
	CentrifugoUrl         string
	CentrifugoApiKey      string
	CentrifugoTokenSecret string
}

type MinIOConfig struct {
	MinioEndpoint        string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	MinioBucketName      string
	MinioUseSSL          bool
}

type BycryptConfig struct {
	BycryptCost   int
	EncryptionKey string
}

type TelegramConfig struct {
	TeleBotToken         string
	TeleChatID           string
	TelegramNotification bool
}

type EmailConfig struct {
	SmtpHost      string
	SmtpPort      string
	SmtpUsername  string
	SmtpPassword  string
	SmtpFromEmail string
	SmtpFromName  string
}

type ScraperConfig struct {
	Host string
	Port string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}
	viper.AutomaticEnv()

	config := &Config{
		App: AppConfig{
			AppName:    getEnv("APP_NAME", "go-service"),
			AppEnv:     getEnv("APP_ENV", "development"),
			IsDebug:    getEnvAsBool("APP_DEBUG", true),
			Port:       getEnv("APP_PORT", "8080"),
			URL:        getEnv("APP_URL", "http://localhost:8080"),
			ClientUrl:  getEnv("APP_CLIENT_URL", "http://localhost:3000"),
			AdminEmail: getEnv("APP_ADMIN_EMAIL", "premiumwatchdevice@gmail.com"),
		},
		Database: DatabaseConfig{
			DbHost:            getEnv("DB_HOST", "localhost"),
			DbPort:            getEnv("DB_PORT", "5432"),
			DbUser:            getEnv("DB_USER", "postgres"),
			DbPassword:        getEnv("DB_PASSWORD", ""),
			DbName:            getEnv("DB_NAME", "indoxxi"),
			DbSSLMode:         getEnv("DB_SSL_MODE", "disable"),
			DbMaxOpenConn:     getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			DbMaxIdleConn:     getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			DbConnMaxLifetime: getAsTime("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			RedisHost:        getEnv("REDIS_HOST", "localhost"),
			RedisPort:        getEnv("REDIS_PORT", "6379"),
			RedisPassword:    getEnv("REDIS_PASSWORD", ""),
			RedisDB:          getEnvAsInt("REDIS_DB", 0),
			RedisAsynqDB:     getEnvAsInt("REDIS_ASYNQ_DB", 1),
			RedisInstance:    getEnv("REDIS_INSTANCE", "default"),
			RedisPoolSize:    getEnvAsInt("REDIS_POOL_SIZE", 10),
			RedisConcurrency: getEnvAsInt("REDIS_CONCURRENCY", 10),
		},
		JWT: JWTConfig{
			JwtSecret:            getEnv("JWT_SECRET", "1234567890"),
			JwtExpiration:        getAsTime("JWT_EXPIRE_HOURS", 15*time.Minute),
			JwtRefreshExpiration: getAsTime("JWT_REFRESH_EXPIRE_HOURS", 24*time.Hour),
		},
		Centrifugo: CentrifugoConfig{
			CentrifugoUrl:         getEnv("CENTRIFUGE_URL", "http://localhost:8000"),
			CentrifugoApiKey:      getEnv("CENTRIFUGO_API_KEY", "1234567890"),
			CentrifugoTokenSecret: getEnv("CENTRIFUGO_TOKEN_SECRET", "1234567890"),
		},
		MinIO: MinIOConfig{
			MinioEndpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			MinioAccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			MinioSecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			MinioBucketName:      getEnv("MINIO_BUCKET", "indoxxi"),
			MinioUseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		},
		Bycrypt: BycryptConfig{
			BycryptCost:   getEnvAsInt("BYCRYPT_COST", 10),
			EncryptionKey: getEnv("ENCRYPTION_KEY", "1234567890"),
		},
		Telegram: TelegramConfig{
			TeleBotToken:         getEnv("TELEGRAM_BOT_TOKEN", "1234567890"),
			TeleChatID:           getEnv("TELEGRAM_CHAT_ID", "-1001234567890"),
			TelegramNotification: getEnvAsBool("TELEGRAM_NOTIFICATIONS", true),
		},
		Email: EmailConfig{
			SmtpHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
			SmtpPort:      getEnv("SMTP_PORT", "587"),
			SmtpUsername:  getEnv("SMTP_USERNAME", ""),
			SmtpPassword:  getEnv("SMTP_PASSWORD", ""),
			SmtpFromEmail: getEnv("SMTP_FROM_EMAIL", ""),
			SmtpFromName:  getEnv("SMTP_FROM_NAME", ""),
		},
		Scraper: ScraperConfig{
			Host: getEnv("SCRAPER_HOST", "localhost"),
			Port: getEnv("SCRAPER_PORT", "50051"),
		},
	}

	return config, nil
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value bool
	fmt.Sscanf(valueStr, "%t", &value)
	return value
}

func getAsTime(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value time.Duration
	fmt.Sscanf(valueStr, "%v", &value)
	return value
}

func getAsInt32(key string, defaultValue int32) int32 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value int32
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

func getAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value int64
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

func (c *DatabaseConfig) GetDbDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSSLMode,
	)
}
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
func (c *AppConfig) IsDevelopment() bool {
	return c.AppEnv == "development"
}
func (c *AppConfig) IsProduction() bool {
	return c.AppEnv == "production"
}
func (c *AppConfig) IsDebugMode() bool {
	return c.IsDebug
}
func (c *RedisConfig) GetRedisInstance() string {
	return c.RedisInstance
}
