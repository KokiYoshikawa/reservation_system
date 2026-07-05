package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"reservation-system/backend/internal/config"
	"reservation-system/backend/internal/domain"
	"reservation-system/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

type authService struct {
	repo      *repository.Repository
	jwtSecret []byte
	jwtIssuer string
	tokenTTL  time.Duration
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Iss   string `json:"iss"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
}

func NewAuthService(repo *repository.Repository, cfg config.Config) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(cfg.JWTSecret),
		jwtIssuer: cfg.JWTIssuer,
		tokenTTL:  cfg.SessionTTL,
	}
}

func (s *authService) Login(ctx context.Context, email string) (string, *domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, ErrUserNotFound
		}

		return "", nil, fmt.Errorf("service login get user by email: %w", err)
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("service login generate jwt: %w", err)
	}

	return token, user, nil
}

func (s *authService) Me(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("service me get user by email: %w", err)
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	claims, err := s.ValidateToken(ctx, token)
	if err != nil {
		return err
	}

	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
	if ttl <= 0 {
		return nil
	}

	if err := s.repo.RevokeToken(ctx, token, ttl); err != nil {
		return fmt.Errorf("service logout revoke token: %w", err)
	}

	return nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*AuthClaims, error) {
	revoked, err := s.repo.IsTokenRevoked(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("service validate token revoked check: %w", err)
	}

	if revoked {
		return nil, ErrTokenRevoked
	}

	payload, err := s.parseJWT(token)
	if err != nil {
		return nil, err
	}

	return &AuthClaims{
		Subject:   payload.Sub,
		Email:     payload.Email,
		IssuedAt:  payload.Iat,
		ExpiresAt: payload.Exp,
	}, nil
}

func (s *authService) generateJWT(user *domain.User) (string, error) {
	now := time.Now().UTC()

	headerPart, err := s.encodeJWTPart(jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	})
	if err != nil {
		return "", err
	}

	payloadPart, err := s.encodeJWTPart(jwtPayload{
		Sub:   fmt.Sprintf("%d", user.ID),
		Email: user.Email,
		Iss:   s.jwtIssuer,
		Iat:   now.Unix(),
		Exp:   now.Add(s.tokenTTL).Unix(),
	})
	if err != nil {
		return "", err
	}

	signingInput := headerPart + "." + payloadPart
	signature := s.signJWT(signingInput)

	return signingInput + "." + signature, nil
}

func (s *authService) parseJWT(token string) (*jwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := s.signJWT(signingInput)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	var payload jwtPayload
	if err := s.decodeJWTPart(parts[1], &payload); err != nil {
		return nil, ErrInvalidToken
	}

	if payload.Iss != s.jwtIssuer || payload.Email == "" || payload.Sub == "" {
		return nil, ErrInvalidToken
	}

	if time.Now().UTC().Unix() >= payload.Exp {
		return nil, ErrInvalidToken
	}

	return &payload, nil
}

func (s *authService) encodeJWTPart(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal jwt part: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (s *authService) decodeJWTPart(encoded string, dest any) error {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, dest)
}

func (s *authService) signJWT(input string) string {
	mac := hmac.New(sha256.New, s.jwtSecret)
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
