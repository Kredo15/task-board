package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *JWTManager) GeneratePair(userID string) (string, string, error) {
	// 1. Создаем Access Token
	accessToken, err := j.generateToken(userID, j.accessTTL)
	if err != nil {
		return "", "", err
	}

	// 2. Создаем Refresh Token
	refreshToken, err := j.generateToken(userID, j.refreshTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTManager) generateToken(userID string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}