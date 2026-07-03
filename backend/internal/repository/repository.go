package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db         *pgxpool.Pool
	redis      *redis.Client
	sessionTTL time.Duration
}

func New(db *pgxpool.Pool, redisClient *redis.Client, sessionTTL time.Duration) *Repository {
	return &Repository{
		db:         db,
		redis:      redisClient,
		sessionTTL: sessionTTL,
	}
}

func (r *Repository) GetRootMessage() string {
	return "reservation system backend is running"
}

func (r *Repository) GetHealthStatus() string {
	return "ok"
}

func (r *Repository) CheckHealth(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r *Repository) CheckRedisHealth(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

func (r *Repository) CreateSession(ctx context.Context, email string) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}

	token := hex.EncodeToString(tokenBytes)
	key := r.sessionKey(token)

	if err := r.redis.Set(ctx, key, email, r.sessionTTL).Err(); err != nil {
		return "", fmt.Errorf("store session: %w", err)
	}

	return token, nil
}

func (r *Repository) GetSessionEmail(ctx context.Context, token string) (string, error) {
	email, err := r.redis.Get(ctx, r.sessionKey(token)).Result()
	if err != nil {
		return "", err
	}

	return email, nil
}

func (r *Repository) DeleteSession(ctx context.Context, token string) error {
	return r.redis.Del(ctx, r.sessionKey(token)).Err()
}

func ExtractBearerToken(headerValue string) string {
	if headerValue == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(headerValue, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(headerValue, prefix))
}

func (r *Repository) sessionKey(token string) string {
	return "auth:session:" + token
}
