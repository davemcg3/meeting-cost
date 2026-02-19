package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription represents an organization's subscription to the service.
type Subscription struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Organization
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_subscription_org" json:"organization_id"`

	// Subscription details
	PlanType           string    `gorm:"type:varchar(50);not null" json:"plan_type"` // "free", "basic", "premium", "enterprise"
	Status             string    `gorm:"type:varchar(50);not null" json:"status"`     // "active", "canceled", "past_due", "trialing"
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`

	// Stripe integration
	StripeCustomerID       string `gorm:"type:varchar(255);uniqueIndex:idx_subscription_stripe_customer" json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID  string `gorm:"type:varchar(255);uniqueIndex:idx_subscription_stripe_sub" json:"stripe_subscription_id,omitempty"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	Payments     []Payment    `gorm:"foreignKey:SubscriptionID" json:"-"`
}

// TableName overrides the table name.
func (Subscription) TableName() string {
	return "subscriptions"
}

// BeforeCreate ensures UUID is set if not already.
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
