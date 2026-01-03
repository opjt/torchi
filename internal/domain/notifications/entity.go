package notifications

import "github.com/google/uuid"

type Noti struct {
	ServiceID uuid.UUID
	Body      string
}
