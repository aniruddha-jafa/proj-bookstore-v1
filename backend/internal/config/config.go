package config

import (
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	AppEnv            string `env:"APP_ENV" env-default:"prod"`
	Port              int    `env:"PORT"`
	JwtSecret         string `env:"JWT_SECRET"`
	MinPasswordLength int    `env:"MIN_PASSWORD_LENGTH"`
	// Comma-separated origins for browser clients (e.g. Next.js dev server).
	CORSAllowOrigin      string        `env:"CORS_ALLOW_ORIGIN"`
	RefreshTokenValidity time.Duration `env:"REFRESH_TOKEN_VALIDITY"`
	AccessTokenValidity  time.Duration `env:"ACCESS_TOKEN_VALIDITY"`
	DbConfig             DbConfig
}

func (appConfig *AppConfig) String() string {
	return fmt.Sprintf("AppEnv: %s, Port: %d, CORSAllowOrigin: %s, RefreshTokenValidity: %s, AccessTokenValidity: %s, DbConfig: %v", appConfig.AppEnv, appConfig.Port, appConfig.CORSAllowOrigin, appConfig.RefreshTokenValidity, appConfig.AccessTokenValidity, appConfig.DbConfig)
}

func (dbConfig *DbConfig) String() string {
	return fmt.Sprintf("DbHost: %s, DbPort: %d, DbUser: %s, DbPassword: %s, DbName: %s, DbSSLMode: %s", dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbName, dbConfig.DbSSLMode)
}

type DbConfig struct {
	DbHost     string `env:"DB_HOST"`
	DbPort     int    `env:"DB_PORT"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
	DbSSLMode  string `env:"DB_SSL_MODE"`
}

var (
	appConfig *AppConfig
	once      sync.Once
)

const MIN_REFRESH_TOKEN_VALIDITY_PROD = time.Minute * 15
const MIN_ACCESS_TOKEN_VALIDITY_PROD = time.Minute * 1

func InitAppConfig() *AppConfig {
	// Initialize it exactly once
	once.Do(func() {
		appConfig = &AppConfig{}
		err := cleanenv.ReadConfig(".env", appConfig)
		if err != nil {
			slog.Error("Unable to read app config", "error", err)
		}
		ValidateSecurityConfig(appConfig)
		slog.Info("App config initialized", "appConfig", appConfig.String())
	})
	return appConfig
}

func ValidateSecurityConfig(appConfig *AppConfig) {
	if appConfig.CORSAllowOrigin == "*" {
		log.Fatal("CORSAllowOrigins cannot be *")
	}
	if appConfig.RefreshTokenValidity <= 0 || appConfig.AccessTokenValidity <= 0 {
		log.Fatal("RefreshTokenValidity and AccessTokenValidity must be set")
	}
	if appConfig.RefreshTokenValidity <= appConfig.AccessTokenValidity {
		log.Fatal("RefreshTokenValidity must be greater than AccessTokenValidity")
	}
	if appConfig.AppEnv == "prod" {
		if appConfig.RefreshTokenValidity < MIN_REFRESH_TOKEN_VALIDITY_PROD {
			appConfig.RefreshTokenValidity = MIN_REFRESH_TOKEN_VALIDITY_PROD
		}
		if appConfig.AccessTokenValidity < MIN_ACCESS_TOKEN_VALIDITY_PROD {
			appConfig.AccessTokenValidity = MIN_ACCESS_TOKEN_VALIDITY_PROD
		}
	}
}
