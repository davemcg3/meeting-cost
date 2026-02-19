package container

import (
	"context"

	"github.com/yourorg/meeting-cost/backend/go/internal/auth"
	"github.com/yourorg/meeting-cost/backend/go/internal/cache"
	"github.com/yourorg/meeting-cost/backend/go/internal/config"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/pubsub"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository/gorm"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
	"github.com/yourorg/meeting-cost/backend/go/internal/service/impl"
	gormio "gorm.io/gorm"
)

// Container manages application dependencies.
type Container struct {
	DB     *gormio.DB
	Cache  cache.Cache
	PubSub pubsub.PubSub
	Logger logger.Logger

	// Repositories
	PersonRepo     repository.PersonRepository
	OrgRepo        repository.OrganizationRepository
	ProfileRepo    repository.PersonOrganizationProfileRepository
	MeetingRepo    repository.MeetingRepository
	IncrementRepo  repository.IncrementRepository
	AuthRepo       repository.AuthRepository
	PermissionRepo repository.PermissionRepository
	ConsentRepo    repository.ConsentRepository
	AuditLogRepo   repository.AuditLogRepository

	// Services
	AuthService     service.AuthService
	PersonService   service.PersonService
	OrgService      service.OrganizationService
	MeetingService  service.MeetingService
	ConsentService  service.ConsentService
	AuditLogService service.AuditLogService
}

// NewContainer initializes all dependencies.
func NewContainer(ctx context.Context, cfg *config.Config, db *gormio.DB, cacheClient cache.Cache, log logger.Logger) (*Container, error) {
	c := &Container{
		DB:     db,
		Cache:  cacheClient,
		Logger: log,
	}

	// Initialize Auth components
	tokenManager := auth.NewTokenManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTIssuer,
		cfg.Auth.AccessExpiry,
		cfg.Auth.RefreshExpiry,
	)

	// Initialize repositories
	c.PersonRepo = gorm.NewPersonRepository(db, cacheClient)
	c.OrgRepo = gorm.NewOrganizationRepository(db, cacheClient)
	c.ProfileRepo = gorm.NewPersonOrganizationProfileRepository(db, cacheClient)
	c.MeetingRepo = gorm.NewMeetingRepository(db, cacheClient)
	c.IncrementRepo = gorm.NewIncrementRepository(db, cacheClient)
	c.AuthRepo = gorm.NewAuthRepository(db, cacheClient)
	c.PermissionRepo = gorm.NewPermissionRepository(db, cacheClient)
	c.ConsentRepo = gorm.NewConsentRepository(db, cacheClient)
	c.AuditLogRepo = gorm.NewAuditLogRepository(db)

	// Initialize PubSub
	c.PubSub = pubsub.NewRedisPubSub(cacheClient.GetClient())

	// Initialize services
	c.AuditLogService = impl.NewAuditLogService(c.AuditLogRepo)
	c.AuthService = impl.NewAuthService(c.PersonRepo, c.AuthRepo, tokenManager, c.AuditLogService, c.Logger)
	c.ConsentService = impl.NewConsentService(c.ConsentRepo, c.AuditLogService)

	c.OrgService = impl.NewOrganizationService(
		c.OrgRepo,
		c.ProfileRepo,
		c.PermissionRepo,
		c.PersonRepo,
		c.AuditLogService,
		c.Logger,
	)

	c.MeetingService = impl.NewMeetingService(
		c.MeetingRepo,
		c.IncrementRepo,
		c.OrgRepo,
		c.ProfileRepo,
		c.PermissionRepo,
		c.AuditLogService,
		c.Cache,
		c.PubSub,
		c.Logger,
	)

	return c, nil
}

// Close performs cleanup of dependencies.
func (c *Container) Close() error {
	// Add cleanup logic if needed (e.g. closing db, cache connections)
	return nil
}
