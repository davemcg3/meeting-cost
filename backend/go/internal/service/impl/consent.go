package impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type consentService struct {
	repo            repository.ConsentRepository
	auditLogService service.AuditLogService
}

func NewConsentService(repo repository.ConsentRepository, auditLogService service.AuditLogService) service.ConsentService {
	return &consentService{
		repo:            repo,
		auditLogService: auditLogService,
	}
}

func (s *consentService) GetConsent(ctx context.Context, sessionID string) (*service.ConsentDTO, error) {
	consent, err := s.repo.GetCurrentBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.mapToDTO(consent), nil
}

func (s *consentService) UpdateConsent(ctx context.Context, req service.UpdateConsentRequest) (*service.ConsentDTO, error) {
	previous, _ := s.repo.GetCurrentBySession(ctx, req.SessionID)

	consent := &models.CookieConsent{
		SessionID:         req.SessionID,
		PersonID:          req.PersonID,
		NecessaryCookies:  true,
		AnalyticsCookies:  req.AnalyticsCookies,
		MarketingCookies:  req.MarketingCookies,
		FunctionalCookies: req.FunctionalCookies,
		ConsentVersion:    "1.0.0", // Hardcoded for now
		ConsentDate:       time.Now(),
		IPAddress:         req.IPAddress,
		UserAgent:         req.UserAgent,
		ConsentSource:     "update",
	}

	if previous != nil {
		consent.PreviousConsentID = &previous.ID
	}

	if err := s.repo.Create(ctx, consent); err != nil {
		return nil, err
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:     consent.PersonID,
		Action:       "update_cookie_consent",
		ResourceType: "cookie_consent",
		ResourceID:   consent.ID,
		IPAddress:    consent.IPAddress,
		UserAgent:    consent.UserAgent,
		Details: map[string]interface{}{
			"analytics":  consent.AnalyticsCookies,
			"marketing":  consent.MarketingCookies,
			"functional": consent.FunctionalCookies,
			"version":    consent.ConsentVersion,
		},
	})

	return s.mapToDTO(consent), nil
}

func (s *consentService) WithdrawConsent(ctx context.Context, sessionID string, cookieTypes []string) error {
	previous, err := s.repo.GetCurrentBySession(ctx, sessionID)
	if err != nil {
		return err
	}

	consent := &models.CookieConsent{
		SessionID:         sessionID,
		PersonID:          previous.PersonID,
		NecessaryCookies:  true,
		AnalyticsCookies:  previous.AnalyticsCookies,
		MarketingCookies:  previous.MarketingCookies,
		FunctionalCookies: previous.FunctionalCookies,
		ConsentVersion:    previous.ConsentVersion,
		ConsentDate:       time.Now(),
		ConsentSource:     "withdrawal",
		PreviousConsentID: &previous.ID,
	}

	for _, ct := range cookieTypes {
		switch ct {
		case "analytics":
			consent.AnalyticsCookies = false
		case "marketing":
			consent.MarketingCookies = false
		case "functional":
			consent.FunctionalCookies = false
		}
	}

	if err := s.repo.Create(ctx, consent); err != nil {
		return err
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:     consent.PersonID,
		Action:       "withdraw_cookie_consent",
		ResourceType: "cookie_consent",
		ResourceID:   consent.ID,
		IPAddress:    consent.IPAddress,
		UserAgent:    consent.UserAgent,
		Details: map[string]interface{}{
			"withdrawn_types": cookieTypes,
		},
	})

	return nil
}

func (s *consentService) CheckCookieAllowed(ctx context.Context, sessionID string, cookieCategory string) (bool, error) {
	if cookieCategory == "necessary" {
		return true, nil
	}

	consent, err := s.repo.GetCurrentBySession(ctx, sessionID)
	if err != nil {
		return false, nil // Default to false if no consent found
	}

	switch cookieCategory {
	case "analytics":
		return consent.AnalyticsCookies, nil
	case "marketing":
		return consent.MarketingCookies, nil
	case "functional":
		return consent.FunctionalCookies, nil
	default:
		return false, nil
	}
}

func (s *consentService) ClassifyCookie(cookieName string) string {
	// Simple classification logic
	switch cookieName {
	case "_ga", "_gid":
		return "analytics"
	case "ads_token":
		return "marketing"
	case "theme", "lang":
		return "functional"
	case "session_id":
		return "necessary"
	default:
		return "necessary"
	}
}

func (s *consentService) GetConsentHistory(ctx context.Context, sessionID string, personID *uuid.UUID) ([]*service.ConsentDTO, error) {
	var models []*models.CookieConsent
	var err error

	if personID != nil {
		models, err = s.repo.GetHistoryByPerson(ctx, *personID)
	} else {
		models, err = s.repo.GetHistoryBySession(ctx, sessionID)
	}

	if err != nil {
		return nil, err
	}

	dtos := make([]*service.ConsentDTO, len(models))
	for i, m := range models {
		dtos[i] = s.mapToDTO(m)
	}
	return dtos, nil
}

func (s *consentService) ExportConsentData(ctx context.Context, personID uuid.UUID) (*service.ConsentExportDTO, error) {
	history, err := s.repo.GetHistoryByPerson(ctx, personID)
	if err != nil {
		return nil, err
	}

	consents := make([]service.ConsentDTO, len(history))
	for i, m := range history {
		consents[i] = *s.mapToDTO(m)
	}

	return &service.ConsentExportDTO{
		PersonID:   personID,
		Consents:   consents,
		ExportDate: time.Now(),
	}, nil
}

func (s *consentService) GetCurrentPolicyVersion(ctx context.Context) (string, error) {
	return "1.0.0", nil
}

func (s *consentService) SyncConsent(ctx context.Context, sessionID string, personID uuid.UUID) error {
	consent, err := s.repo.GetCurrentBySession(ctx, sessionID)
	if err != nil {
		return nil // No consent to sync
	}

	// If already tied to this person, nothing to do
	if consent.PersonID != nil && *consent.PersonID == personID {
		return nil
	}

	// Create a new consent record tied to the person
	newConsent := &models.CookieConsent{
		SessionID:         sessionID,
		PersonID:          &personID,
		NecessaryCookies:  true,
		AnalyticsCookies:  consent.AnalyticsCookies,
		MarketingCookies:  consent.MarketingCookies,
		FunctionalCookies: consent.FunctionalCookies,
		ConsentVersion:    consent.ConsentVersion,
		ConsentDate:       time.Now(),
		ConsentSource:     "sync",
		PreviousConsentID: &consent.ID,
	}

	if err := s.repo.Create(ctx, newConsent); err != nil {
		return err
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:     newConsent.PersonID,
		Action:       "sync_cookie_consent",
		ResourceType: "cookie_consent",
		ResourceID:   newConsent.ID,
		Details: map[string]interface{}{
			"session_id": sessionID,
		},
	})

	return nil
}

func (s *consentService) mapToDTO(m *models.CookieConsent) *service.ConsentDTO {
	return &service.ConsentDTO{
		ID:                m.ID,
		PersonID:          m.PersonID,
		SessionID:         m.SessionID,
		NecessaryCookies:  m.NecessaryCookies,
		AnalyticsCookies:  m.AnalyticsCookies,
		MarketingCookies:  m.MarketingCookies,
		FunctionalCookies: m.FunctionalCookies,
		ConsentVersion:    m.ConsentVersion,
		ConsentDate:       m.ConsentDate,
		PreviousConsentID: m.PreviousConsentID,
	}
}
