package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type MeetingHandler struct {
	meetingService service.MeetingService
}

func NewMeetingHandler(s service.MeetingService) *MeetingHandler {
	return &MeetingHandler{
		meetingService: s,
	}
}

func (h *MeetingHandler) CreateMeeting(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)

	var req service.CreateMeetingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	meeting, err := h.meetingService.CreateMeeting(c.Context(), req.OrganizationID, personID, req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(meeting)
}

func (h *MeetingHandler) GetMeeting(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	meeting, err := h.meetingService.GetMeeting(c.Context(), id, personID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(meeting)
}

func (h *MeetingHandler) StartMeeting(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	if err := h.meetingService.StartMeeting(c.Context(), id, personID); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MeetingHandler) StopMeeting(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	if err := h.meetingService.StopMeeting(c.Context(), id, personID); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MeetingHandler) UpdateAttendeeCount(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	var req struct {
		Count int `json:"count"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.meetingService.UpdateAttendeeCount(c.Context(), id, req.Count, personID, c.IP(), string(c.Request().Header.UserAgent())); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MeetingHandler) GetMeetingCost(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	res, err := h.meetingService.GetMeetingCost(c.Context(), id, personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *MeetingHandler) ListMeetings(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)

	orgIDStr := c.Query("organization_id")
	if orgIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organization_id is required"})
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization_id"})
	}

	filters := service.MeetingFilters{}
	pagination := service.Pagination{Page: 1, PageSize: 100}

	res, _, err := h.meetingService.ListMeetings(c.Context(), orgID, personID, filters, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
func (h *MeetingHandler) DeleteMeeting(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	if err := h.meetingService.DeleteMeeting(c.Context(), id, personID, c.IP(), string(c.Request().Header.UserAgent())); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
