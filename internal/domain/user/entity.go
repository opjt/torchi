package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       *string
	IsGuest     bool
	TermsAgreed bool
	CreatedAt   time.Time
}
