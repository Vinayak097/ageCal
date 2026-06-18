package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/dob-api/internal/handler"
	"github.com/user/dob-api/internal/middleware"
	"go.uber.org/zap"
)

func Register(app *fiber.App, h *handler.UserHandler, log *zap.Logger) {
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(log))

	v1 := app.Group("/users")
	v1.Post("/", h.CreateUser)
	v1.Get("/", h.ListUsers)
	v1.Get("/:id", h.GetUser)
	v1.Put("/:id", h.UpdateUser)
	v1.Delete("/:id", h.DeleteUser)
}
