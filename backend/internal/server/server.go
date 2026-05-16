package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/auth"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/middleware"
	"github.com/aniruddha-jafa/go-auth-v1/internal/refresh_tokens"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // or your preferred SQL driver
)

func errorHandler(c *fiber.Ctx, err error) error {
	log.Printf("ERROR: %v", err)

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
	pool, err := initDbConfig(config.DbConfig)
	if err != nil {
		panic(err)
	}
	queries := db.New(pool)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

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
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		err := ctx.SendString("pong\n")
		return err
	})

	// User
	userRepo := users.NewUserRepoImpl(queries)
	userService := users.UserServiceImpl{
		UserRepo: &userRepo,
	}
	userHandler := users.UserHandlerImpl{
		UserService: &userService,
	}
	// Auth
	refreshTokenRepo := refresh_tokens.NewRefreshTokenRepoImpl(queries)
	authService := auth.AuthServiceImpl{
		UserService:      &userService,
		RefreshTokenRepo: &refreshTokenRepo,
	}
	authHandler := auth.AuthHandlerImpl{
		AuthService: &authService,
	}

	// Routes
	apiGroup := app.Group("api")
	apiV1Group := apiGroup.Group("v1")

	// Auth routes
	authGroup := apiV1Group.Group("auth")
	authGroup.Post("/signup", authHandler.SignUp)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh-token", authHandler.RefreshToken)
	authGroup.Post("/logout", authHandler.Logout)

	// User routes
	userGroup := apiV1Group.Group("users")
	userGroup.Get("/:id", middleware.RequireAuth, userHandler.Get)
	userGroup.Post("/:id", middleware.RequireAuth, userHandler.Update)
	userGroup.Get("/", userHandler.GetAll)

	port := ":" + strconv.Itoa(config.Port)
	log.Printf("Listening on port: %s", port)
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
	log.Printf("Connecting to DB with connection string: %s", connStr)

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
