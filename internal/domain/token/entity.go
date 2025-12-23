package token

import "github.com/google/uuid"

type Token struct {
	P256dh string
	Auth   string
	UserID uuid.UUID
}
