package database

import (
	"time"

	"emperror.dev/errors"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	CreatedAt time.Time  `gorm:"index"       json:"createdAt"`
	UpdatedAt time.Time  `gorm:"index"       json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index"       json:"deletedAt,omitempty"`
	ID        string     `gorm:"primary_key" json:"id"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(_ *gorm.DB) error {
	// Check if ID is set to avoid erasing it.
	// This is useful when it is asked to save object for the first
	// time with a fixed id.
	if base.ID == "" {
		// Generate new id
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return errors.WithStack(err)
		}

		// Save new id
		base.ID = uuidGenerated.String()
	}

	return nil
}
