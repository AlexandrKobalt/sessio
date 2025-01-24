package models

import (
	"encoding/json"
	"time"
)

type Session struct {
	UserID    string    `json:"user_id"`
	AuthKey   string    `json:"auth_key"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, s)
}
