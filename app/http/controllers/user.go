package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/views"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// UserController is the controller for the user routes
type UserController struct {
	user.IUseCase
}

// NewUserController creates a new user controller
func NewUserController(useCase user.IUseCase) UserController {
	return UserController{IUseCase: useCase}
}

// GetName returns the name of the user
func (u *UserController) GetName(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	return c.SendString(u.IUseCase.GetName(c.UserContext(), uuid))
}

// GetMe returns the user from the context
func (*UserController) GetMe(c *fiber.Ctx) error {
	userCtx, err := GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVM("User not found", nil))
	}

	otelzap.L().Info("user found", zap.String("user", userCtx.Username))
	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("User found", userCtx))
}
