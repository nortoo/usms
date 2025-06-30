package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/log"
	"github.com/nortoo/usms/internal/pkg/session"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/nortoo/usms/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Claims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for access and refresh token.
func GenerateToken(tokenId, secret string, userID uint, expiryIn int64) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryIn) * time.Second) // Access token valid for 15 minutes

	claims := &Claims{
		UID: int64(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			ID:        tokenId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses and validates an jwt token
func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}
	return claims, nil
}

func IssueAccessTokenAndRefreshToken(uid uint) (accessToken, refreshToken string, err error) {
	tokenId := uuid.NewV4().String()

	accessToken, err = GenerateToken(
		tokenId,
		etc.GetEnv().JWTSecretKey,
		uid,
		etc.GetConfig().App.Settings.JWT.TokenExpireTime)
	if err != nil {
		return "", "", errors.ErrInternalError.WithDetail(err.Error())
	}

	refreshToken, err = GenerateToken(
		tokenId,
		etc.GetEnv().JWTRefreshSecretKey,
		uid,
		etc.GetConfig().App.Settings.JWT.TokenRefreshTime)
	if err != nil {
		return "", "", errors.ErrInternalError.WithDetail(err.Error())
	}

	sessionStoreKey := session.GenerateSessionStoreKey(uid, tokenId)
	if err = store.GetRedisClient().GetRDB().Set(
		context.TODO(),
		sessionStoreKey,
		"",
		time.Duration(etc.GetConfig().App.Settings.JWT.TokenExpireTime)*time.Second).Err(); err != nil {
		log.GetLogger().Warn("failed to store session.", zap.String("key", sessionStoreKey), zap.Error(err))
	}

	refreshTokenStoreKey := session.GenerateSessionRefreshTokenStoreKey(uid, tokenId)
	if err = store.GetRedisClient().GetRDB().Set(
		context.TODO(),
		refreshTokenStoreKey,
		"",
		time.Duration(etc.GetConfig().App.Settings.JWT.TokenRefreshTime)*time.Second).Err(); err != nil {
		log.GetLogger().Warn("failed to store refresh token.", zap.String("key", sessionStoreKey), zap.Error(err))
	}

	return
}
