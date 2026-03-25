package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       *string
	Provider    *string
	ProviderID  *string
	IsGuest     bool
	TermsAgreed bool
	CreatedAt   time.Time
}
