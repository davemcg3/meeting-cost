package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment represents a payment transaction.
type Payment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Subscription
	SubscriptionID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_subscription" json:"subscription_id"`

	// Payment details
	Amount   float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency string     `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Status   string     `gorm:"type:varchar(50);not null" json:"status"` // "succeeded", "pending", "failed", "refunded"
	PaidAt   *time.Time `json:"paid_at,omitempty"`

	// Stripe integration
	StripePaymentIntentID string `gorm:"type:varchar(255);uniqueIndex:idx_payment_stripe" json:"stripe_payment_intent_id,omitempty"`
	ReceiptURL            string `gorm:"type:text" json:"receipt_url,omitempty"`

	// Relationships
	Subscription Subscription `gorm:"foreignKey:SubscriptionID" json:"-"`
}

// TableName overrides the table name.
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate ensures UUID is set if not already.
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
