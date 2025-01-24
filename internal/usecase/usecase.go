package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/AlexandrKobalt/sessio/internal"
	"github.com/AlexandrKobalt/sessio/internal/models"
	"github.com/google/uuid"
)

var (
	ErrTooManySessions  = errors.New("too many sessions")
	ErrInvalidUserAgent = errors.New("invalid user agent")
)

type useCase struct {
	maxSessionsCount int
	cache            internal.Cache
}

func New(maxSessionsCount int, cache internal.Cache) internal.UseCase {
	return &useCase{maxSessionsCount: maxSessionsCount, cache: cache}
}

func (uc *useCase) CreateSession(ctx context.Context, userID, userAgent string) (authKey string, err error) {
	sessions, err := uc.cache.GetAllSessions(ctx, userID)
	if err != nil {
		return "", err
	}

	for _, s := range sessions {
		if s.UserAgent == userAgent {
			return s.AuthKey, nil
		}
	}

	if len(sessions) >= uc.maxSessionsCount {
		return "", ErrTooManySessions
	}

	authKey = uuid.NewString()
	now := time.Now()

	if err = uc.cache.CreateSession(
		ctx,
		models.Session{
			UserID:    userID,
			UserAgent: userAgent,
			AuthKey:   authKey,
			CreatedAt: now,
			ExpiresAt: now.Add(time.Hour * 24 * 30),
		},
	); err != nil {
		return "", nil
	}

	return authKey, nil
}

func (uc *useCase) Authenticate(ctx context.Context, authKey, userAgent string) (userID string, err error) {
	session, err := uc.cache.GetSession(ctx, authKey)
	if err != nil {
		return "", err
	}

	if session.UserAgent != userAgent {
		return "", ErrInvalidUserAgent
	}

	return session.UserID, nil
}

func (uc *useCase) GetAllActiveSessions(ctx context.Context, userID string) ([]models.Session, error) {
	return uc.cache.GetAllSessions(ctx, userID)
}

func (uc *useCase) DeleteSession(ctx context.Context, userID string, userAgent string) (err error) {
	return uc.cache.DeleteSession(ctx, userID, userAgent)
}
