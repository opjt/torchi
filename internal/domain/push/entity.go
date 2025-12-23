package push

import "github.com/google/uuid"

type Subscription struct {
	UserID   uuid.UUID
	Endpoint string
	P256dh   string
	Auth     string
}
