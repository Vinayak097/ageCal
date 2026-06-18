package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/user/dob-api/internal/models"
	"github.com/user/dob-api/internal/repository"
	"github.com/user/dob-api/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
	log      *zap.Logger
}

func NewUserHandler(svc service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:      svc,
		validate: validator.New(),
		log:      log,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.svc.CreateUser(c.Context(), req)
	if err != nil {
		h.log.Error("CreateUser failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	h.log.Info("user created", zap.Int32("id", resp.ID))
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	resp, err := h.svc.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
		}
		h.log.Error("GetUser failed", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	return c.JSON(resp)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
		}
		h.log.Error("UpdateUser failed", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	h.log.Info("user updated", zap.Int32("id", resp.ID))
	return c.JSON(resp)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid id"})
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		h.log.Error("DeleteUser failed", zap.Error(err), zap.Int32("id", id))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	h.log.Info("user deleted", zap.Int32("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page := int32(c.QueryInt("page", 1))
	limit := int32(c.QueryInt("limit", 10))

	resp, err := h.svc.ListUsers(c.Context(), page, limit)
	if err != nil {
		h.log.Error("ListUsers failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
	}

	return c.JSON(resp)
}

func parseID(c *fiber.Ctx) (int32, error) {
	raw := c.Params("id")
	id, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return int32(id), nil
}
