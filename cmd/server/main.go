package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/user/dob-api/config"
	db "github.com/user/dob-api/db/sqlc"
	"github.com/user/dob-api/internal/handler"
	"github.com/user/dob-api/internal/logger"
	"github.com/user/dob-api/internal/repository"
	"github.com/user/dob-api/internal/routes"
	"github.com/user/dob-api/internal/service"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	logger.Log.Info("connecting to database",
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("user", cfg.DBUser),
		zap.String("db", cfg.DBName),
	)
	sqlDB, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		logger.Log.Fatal("failed to open database", zap.Error(err))
	}
	defer sqlDB.Close()

	if err := sqlDB.PingContext(context.Background()); err != nil {
		logger.Log.Fatal("database ping failed", zap.Error(err))
	}
	logger.Log.Info("database connected")

	queries := db.New(sqlDB)
	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, logger.Log)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if ok := errors.As(err, &e); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Register(app, userHandler, logger.Log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		logger.Log.Info("server starting", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			logger.Log.Error("server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Log.Info("shutting down server")
	if err := app.Shutdown(); err != nil {
		logger.Log.Error("shutdown error", zap.Error(err))
	}
}
