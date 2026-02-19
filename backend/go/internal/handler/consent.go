package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type ConsentHandler struct {
	service service.ConsentService
}

func NewConsentHandler(service service.ConsentService) *ConsentHandler {
	return &ConsentHandler{
		service: service,
	}
}

func (h *ConsentHandler) GetConsent(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session_id is required"})
	}

	consent, err := h.service.GetConsent(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "consent not found"})
	}

	return c.JSON(consent)
}

func (h *ConsentHandler) UpdateConsent(c *fiber.Ctx) error {
	var req service.UpdateConsentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Enrich request with context info
	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	// If there's a person in locals (from auth middleware), set it
	if personIDStr, ok := c.Locals("personID").(string); ok {
		if id, err := uuid.Parse(personIDStr); err == nil {
			req.PersonID = &id
		}
	}

	consent, err := h.service.UpdateConsent(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(consent)
}

func (h *ConsentHandler) GetHistory(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")

	var personID *uuid.UUID
	if personIDStr, ok := c.Locals("personID").(string); ok {
		if id, err := uuid.Parse(personIDStr); err == nil {
			personID = &id
		}
	}

	if sessionID == "" && personID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "sessionID or authenticated user required"})
	}

	history, err := h.service.GetConsentHistory(c.Context(), sessionID, personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(history)
}
func (h *ConsentHandler) SyncConsent(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "session_id is required"})
	}

	personIDStr, ok := c.Locals("personID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person_id"})
	}

	if err := h.service.SyncConsent(c.Context(), sessionID, personID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
