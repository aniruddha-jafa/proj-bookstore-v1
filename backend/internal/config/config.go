package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Port              int    `env:"PORT"`
	JwtSecret         string `env:"JWT_SECRET"`
	MinPasswordLength int    `env:"MIN_PASSWORD_LENGTH"`
	// Comma-separated origins for browser clients (e.g. Next.js dev server).
	CORSAllowOrigins string `env:"CORS_ALLOW_ORIGINS"`
	DbConfig         DbConfig
}

func (appConfig *AppConfig) String() string {
	return fmt.Sprintf("Port: %d, DbConfig: %v", appConfig.Port, appConfig.DbConfig)
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

func InitAppConfig() *AppConfig {
	// Initialize it exactly once
	once.Do(func() {
		appConfig = &AppConfig{}
		err := cleanenv.ReadConfig(".env", appConfig)
		if err != nil {
			log.Fatalf("Unable to read app config: %v", err)
		}
		log.Printf("App config initialized: %s", appConfig)
	})
	return appConfig
}
