package endpoint

import (
	"time"

	"github.com/google/uuid"
)

type Endpoint struct {
	ID        uuid.UUID
	Name      string
	Endpoint  string
	CreatedAt time.Time
}
