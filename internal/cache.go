package internal

import (
	"context"

	"github.com/AlexandrKobalt/sessio/internal/models"
)

type Cache interface {
	CreateSession(ctx context.Context, session models.Session) (err error)
	GetAllSessions(ctx context.Context, userID string) (sessions []models.Session, err error)
	GetSession(ctx context.Context, authKey string) (session models.Session, err error)
	DeleteSession(ctx context.Context, userID, userAgent string) (err error)
}
