package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/request_context"
)

const LoggerNameKey = "name"

func InitDefaultLogger(appConfig *config.AppConfig) {
	if appConfig.AppEnv == "prod" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	slog.Info("Logger initialized", "appEnv", appConfig.AppEnv)
}

// Add context values to the logger
func WithContext(baseLogger *slog.Logger, ctx context.Context) *slog.Logger {
	logger := baseLogger

	if requestId, ok := ctx.Value(request_context.RequestIdKey).(string); ok {
		logger = logger.With(string(request_context.RequestIdKey), requestId)
	}

	return logger
}
