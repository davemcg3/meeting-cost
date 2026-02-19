package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	res, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		// In a real app, distinguish between invalid creds and server errors
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(res)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := ""
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	}

	if token != "" {
		_ = h.authService.Logout(c.Context(), token, c.IP(), string(c.Request().Header.UserAgent()))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing refresh token"})
	}

	res, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	return c.JSON(res)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	personID := c.Locals("person_id")
	if personID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// For now just return IDs. Normally we'd call PersonService.
	return c.JSON(fiber.Map{
		"person_id": personID,
		"email":     c.Locals("email"),
	})
}
