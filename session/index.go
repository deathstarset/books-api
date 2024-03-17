package session

import (
	"github.com/google/uuid"
)

type Session struct {
	UserID uuid.UUID `json:"user_id"`
}
