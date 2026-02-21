package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Person represents a user/person in the system. Can belong to multiple organizations.
type Person struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Identity fields
	Email     string `gorm:"uniqueIndex:idx_person_email;not null" json:"email"`
	FirstName string `gorm:"not null" json:"firstName"`
	LastName  string `json:"lastName"`

	// GDPR compliance
	AnonymizedAt *time.Time `json:"anonymized_at,omitempty"` // Set when person requests data deletion
	Anonymized   bool       `gorm:"default:false;index:idx_person_anonymized" json:"anonymized"`

	// Metadata
	Timezone string `gorm:"default:'UTC'" json:"timezone"`
	Locale   string `gorm:"default:'en-US'" json:"locale"`
}

// TableName overrides the table name.
func (Person) TableName() string {
	return "persons"
}

// BeforeCreate ensures UUID is set if not already.
func (p *Person) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
