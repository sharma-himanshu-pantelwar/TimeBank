package session

import (
	"time"

	"github.com/google/uuid"
)

// create a session struct
type Session struct {
	Id        uuid.UUID `json:"id"`
	Uid       int       `json:"uid"`
	TokenHash string    `json:"tokenhash"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}
