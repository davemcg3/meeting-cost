package config

import (
	"fmt"

	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a GORM database connection with connection pooling.
func NewDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// AutoMigrate runs GORM AutoMigrate for all models (development only).
// Production should use versioned SQL migrations.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Person{},
		&models.Organization{},
		&models.PersonOrganizationProfile{},
		&models.Role{},
		&models.RoleAssignment{},
		&models.Permission{},
		&models.AuthMethod{},
		&models.Session{},
		&models.Subscription{},
		&models.Payment{},
		&models.Meeting{},
		&models.Increment{},
		&models.MeetingParticipant{},
		&models.AuditLog{},
		&models.CookieConsent{},
	)
}
