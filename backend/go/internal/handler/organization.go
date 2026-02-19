package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type OrganizationHandler struct {
	orgService service.OrganizationService
}

func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)

	var req service.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	res, err := h.orgService.CreateOrganization(c.Context(), personID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}

	res, err := h.orgService.GetOrganization(c.Context(), orgID, personID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *OrganizationHandler) ListOrganizations(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)

	res, err := h.orgService.ListOrganizations(c.Context(), personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *OrganizationHandler) UpdateOrganization(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}

	var req service.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	res, err := h.orgService.UpdateOrganization(c.Context(), orgID, personID, req)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *OrganizationHandler) GetMembers(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}

	res, err := h.orgService.GetMembers(c.Context(), orgID, personID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *OrganizationHandler) AddMember(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}

	var req service.AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.IPAddress = c.IP()
	req.UserAgent = string(c.Request().Header.UserAgent())

	err = h.orgService.AddMember(c.Context(), orgID, personID, req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid member id"})
	}

	err = h.orgService.RemoveMember(c.Context(), orgID, personID, memberID, c.IP(), string(c.Request().Header.UserAgent()))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *OrganizationHandler) UpdateMemberWage(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}
	memberID, err := uuid.Parse(c.Params("memberId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid member id"})
	}

	var req struct {
		Wage float64 `json:"wage"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err = h.orgService.UpdateMemberWage(c.Context(), orgID, memberID, req.Wage, personID, c.IP(), string(c.Request().Header.UserAgent()))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *OrganizationHandler) DeleteOrganization(c *fiber.Ctx) error {
	personID := c.Locals("person_id").(uuid.UUID)
	orgID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization id"})
	}

	err = h.orgService.DeleteOrganization(c.Context(), orgID, personID, c.IP(), string(c.Request().Header.UserAgent()))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
