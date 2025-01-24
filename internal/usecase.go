package internal

import (
	"context"

	"github.com/AlexandrKobalt/sessio/internal/models"
)

type UseCase interface {
	CreateSession(ctx context.Context, userID, userAgent string) (authKey string, err error)
	Authenticate(ctx context.Context, authKey, userAgent string) (userID string, err error)
	GetAllActiveSessions(ctx context.Context, userID string) ([]models.Session, error)
	DeleteSession(ctx context.Context, userID string, userAgent string) (err error)
}
