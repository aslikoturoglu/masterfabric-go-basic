package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/masterfabric/masterfabric_go_basic/internal/shared/config"
	domainErr "github.com/masterfabric/masterfabric_go_basic/internal/shared/errors"
	"github.com/redis/go-redis/v9"
)

// TokenPair holds a short-lived access token and long-lived refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expiry
}

// Claims is the JWT claims payload.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService handles JWT creation and validation plus refresh token storage in Redis.
type JWTService struct {
	cfg   config.JWTConfig
	redis *redis.Client
}

// NewJWTService creates a JWTService. redis may be nil (tokens won't be blacklisted).
func NewJWTService(cfg config.JWTConfig, redisClient *redis.Client) *JWTService {
	return &JWTService{cfg: cfg, redis: redisClient}
}

// GenerateTokenPair creates an access + refresh token pair for the given user.
func (s *JWTService) GenerateTokenPair(ctx context.Context, userID uuid.UUID, email, role string) (*TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := Claims{
		UserID: userID.String(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Refresh token (opaque UUID stored in Redis as "email:role")
	refreshToken := uuid.New().String()
	if s.redis != nil {
		key := refreshKey(userID.String(), refreshToken)
		value := email + ":" + role
		if err := s.redis.Set(ctx, key, value, s.cfg.RefreshTokenTTL).Err(); err != nil {
			return nil, fmt.Errorf("store refresh token: %w", err)
		}
	}

	return &TokenPair{
		AccessToken:  accessSigned,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}

// ValidateAccessToken parses and validates a JWT access token.
func (s *JWTService) ValidateAccessToken(ctx context.Context, tokenStr string) (*Claims, error) {
	// Check blacklist first
	if s.redis != nil {
		blacklisted, _ := s.redis.Exists(ctx, blacklistKey(tokenStr)).Result()
		if blacklisted > 0 {
			return nil, domainErr.ErrTokenInvalid
		}
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domainErr.ErrTokenExpired
		}
		return nil, domainErr.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domainErr.ErrTokenInvalid
	}
	return claims, nil
}

// RefreshTokens validates a refresh token and issues a new pair.
func (s *JWTService) RefreshTokens(ctx context.Context, userID uuid.UUID, refreshToken string) (*TokenPair, error) {
	if s.redis == nil {
		return nil, domainErr.ErrTokenInvalid
	}

	key := refreshKey(userID.String(), refreshToken)
	value, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, domainErr.ErrTokenInvalid
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	// value is "email:role"
	parts := strings.SplitN(value, ":", 2)
	email := parts[0]
	role := "user"
	if len(parts) == 2 {
		role = parts[1]
	}

	// Rotate: delete old refresh token
	_ = s.redis.Del(ctx, key)

	return s.GenerateTokenPair(ctx, userID, email, role)
}

// RevokeTokens blacklists the access token and deletes the refresh token.
func (s *JWTService) RevokeTokens(ctx context.Context, userID uuid.UUID, accessToken, refreshToken string) error {
	if s.redis == nil {
		return nil
	}

	// Parse to get remaining TTL for blacklist expiry
	claims, err := s.ValidateAccessToken(ctx, accessToken)
	if err == nil && claims != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			_ = s.redis.Set(ctx, blacklistKey(accessToken), "1", ttl)
		}
	}

	// Delete refresh token
	_ = s.redis.Del(ctx, refreshKey(userID.String(), refreshToken))
	return nil
}

func refreshKey(userID, token string) string {
	return fmt.Sprintf("mf:refresh:%s:%s", userID, token)
}

func blacklistKey(token string) string {
	return fmt.Sprintf("mf:blacklist:%s", token)
}
