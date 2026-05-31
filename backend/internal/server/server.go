package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/auth"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/constants"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/internal/middleware"
	"github.com/aniruddha-jafa/go-auth-v1/internal/refresh_tokens"
	"github.com/aniruddha-jafa/go-auth-v1/internal/request_context"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // or your preferred SQL driver
)

func errorHandler(c *fiber.Ctx, err error) error {
	slog.Error("Error", "error", err)

	// Check if it's an HttpError
	if httpErr, ok := err.(apperrors.HttpError); ok {
		return c.Status(httpErr.Status).JSON(apperrors.ErrorResponse{
			Message: httpErr.Message,
			Code:    httpErr.Status,
		})
	}

	var httpErr apperrors.HttpError
	if errors.As(err, &httpErr) {
		return c.Status(httpErr.Status).JSON(apperrors.ErrorResponse{
			Message: httpErr.Message,
			Code:    httpErr.Status,
		})
	}
	var errMessage = "Internal server error"
	if err.Error() != "" {
		errMessage = err.Error()
	}
	// Default to 500 for unknown errors
	return c.Status(http.StatusInternalServerError).JSON(apperrors.ErrorResponse{
		Message: errMessage,
		Code:    http.StatusInternalServerError,
	})
}

func InitServer() {
	config := config.InitAppConfig()
	logger.InitDefaultLogger(config)

	pool, err := initDbConfig(config.DbConfig)
	if err != nil {
		panic(err)
	}
	queries := db.New(pool)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(middleware.RequestID)

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CORSAllowOrigin,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			constants.CSRF_TOKEN_HEADER,
		}, ","),
		AllowCredentials: true,
	}))

	app.Use(loggerMiddleware.New(
		loggerMiddleware.Config{
			CustomTags: map[string]loggerMiddleware.LogFunc{
				string(request_context.RequestIdKey): func(output loggerMiddleware.Buffer, c *fiber.Ctx, data *loggerMiddleware.Data, extraParam string) (int, error) {
					return output.WriteString(c.Get(constants.REQUEST_ID_HEADER, ""))
				},
			},
			Format: "{time: ${time}, method: ${method}, path: ${path}, status: ${status}, requestId: ${requestId}, ip: ${ip}, latency: ${latency}, error: ${error}}",
		},
	))

	// Ping
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		err := ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "pong", "timestamp": time.Now().Format(time.RFC3339)})
		return err
	})

	// User
	userRepo := users.NewUserRepoImpl(queries)
	userService := users.NewUserServiceImpl(userRepo)
	userHandler := users.NewUserHandlerImpl(userService)

	// Auth
	refreshTokenRepo := refresh_tokens.NewRefreshTokenRepoImpl(queries)
	authService := auth.NewAuthServiceImpl(userService, refreshTokenRepo)
	authHandler := auth.NewAuthHandlerImpl(authService)

	// Routes
	apiGroup := app.Group(constants.API)
	apiV1Group := apiGroup.Group(constants.V1)

	// Auth routes
	authGroup := apiV1Group.Group(constants.AUTH)

	// Don't cache on auth routes
	authGroup.Use(middleware.NoCache)

	authGroup.Post(constants.SIGNUP, authHandler.SignUp)
	authGroup.Post(constants.LOGIN, authHandler.Login)
	authGroup.Post(constants.REFRESH_TOKEN, middleware.CSRF, authHandler.RefreshToken)
	authGroup.Post(constants.LOGOUT, middleware.CSRF, authHandler.Logout)
	authGroup.Get(constants.CSRF_TOKEN, authHandler.GetOrResetCSRFToken)

	// User routes
	userGroup := apiV1Group.Group(constants.USER)
	userGroup.Get(":id", middleware.RequireAuth, userHandler.Get)
	userGroup.Post(":id", middleware.RequireAuth, userHandler.Update)
	userGroup.Get("/", userHandler.GetAll)

	port := ":" + strconv.Itoa(config.Port)
	slog.Info("Listening on port", "port", port)
	log.Fatal(app.Listen(port))
}

func initDbConfig(dbConfig config.DbConfig) (*pgxpool.Pool, error) {
	connStr := ""
	if dbConfig.DbHost != "" {
		connStr += fmt.Sprintf("host=%s ", dbConfig.DbHost)
	}
	if dbConfig.DbName != "" {
		connStr += fmt.Sprintf("dbname=%s ", dbConfig.DbName)
	}
	if dbConfig.DbUser != "" {
		connStr += fmt.Sprintf("user=%s ", dbConfig.DbUser)
	}
	if dbConfig.DbPassword != "" {
		connStr += fmt.Sprintf("password=%s ", dbConfig.DbPassword)
	}
	if dbConfig.DbSSLMode != "" {
		connStr += fmt.Sprintf("sslmode=%s ", dbConfig.DbSSLMode)
	}
	connStr = strings.TrimSpace(connStr)
	slog.Info("Connecting to DB", "connectionString", connStr)

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
