package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlexandrKobalt/sessio/internal"
	"github.com/AlexandrKobalt/sessio/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	sessionKeyF     = "session:%s"
	userSessionKeyF = "user:%s:sessions"
)

func createSessionKey(authKey string) string {
	return fmt.Sprintf(sessionKeyF, authKey)
}

func createUserSessionsKey(userID string) string {
	return fmt.Sprintf(userSessionKeyF, userID)
}

type cache struct {
	client *redis.Client
}

func New(client *redis.Client) internal.Cache {
	return &cache{client: client}
}

func (c *cache) CreateSession(ctx context.Context, session models.Session) (err error) {
	sessionKey := createSessionKey(session.AuthKey)
	userSessionsKey := createUserSessionsKey(session.UserID)

	if err = c.client.Set(
		ctx,
		sessionKey,
		session,
		session.ExpiresAt.Sub(session.CreatedAt),
	).Err(); err != nil {
		return err
	}

	return c.client.SAdd(ctx, userSessionsKey, session.AuthKey).Err()
}

func (c *cache) GetAllSessions(ctx context.Context, userID string) (sessions []models.Session, err error) {
	userSessionKey := createUserSessionsKey(userID)

	authKeys, err := c.client.SMembers(ctx, userSessionKey).Result()
	if err != nil {
		return nil, err
	}

	for _, authKey := range authKeys {
		sessionKey := createSessionKey(authKey)

		data, err := c.client.Get(ctx, sessionKey).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var session models.Session
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (c *cache) GetSession(ctx context.Context, authKey string) (session models.Session, err error) {
	sessionKey := createSessionKey(authKey)

	data, err := c.client.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		return session, nil
	} else if err != nil {
		return session, err
	}

	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return session, err
	}

	return session, nil
}

func (c *cache) DeleteSession(ctx context.Context, userID, userAgent string) (err error) {
	userSessionsKey := createUserSessionsKey(userID)

	authKeys, err := c.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return err
	}

	for _, authKey := range authKeys {
		sessionKey := createSessionKey(authKey)
		data, err := c.client.Get(ctx, sessionKey).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return err
		}

		var session models.Session
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			return err
		}

		if session.UserAgent == userAgent {
			if err := c.client.Del(ctx, sessionKey).Err(); err != nil {
				return err
			}

			if err := c.client.SRem(ctx, userSessionsKey, authKey).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}
