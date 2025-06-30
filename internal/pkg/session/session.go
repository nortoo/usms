package session

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/log"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	sessionStoreKeyPrefix             = "session:user"
	sessionRefreshTokenStoreKeyPrefix = "session:refresh"

	sessionBlocklistKeyPrefix             = "session:blocklist:token"
	sessionRefreshTokenBlocklistKeyPrefix = "session:blocklist:refresh"
)

func GenerateSessionStoreKey(uid uint, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", sessionStoreKeyPrefix, uid, tokenId)
}

func GenerateSessionBlocklistKey(tokenId string) string {
	return fmt.Sprintf("%s:%s", sessionBlocklistKeyPrefix, tokenId)
}

func GetTokenIdFromSessionKey(tokenKey string) (string, error) {
	if !strings.HasPrefix(tokenKey, sessionStoreKeyPrefix) {
		return "", errors.New("invalid session key")
	}

	elements := strings.Split(tokenKey, ":")
	if len(elements) != 4 {
		return "", errors.New("invalid session key")
	}

	tokenId := elements[3]
	if tokenId == "" {
		return "", errors.New("invalid session key")
	}
	return tokenId, nil
}

func GenerateSessionRefreshTokenStoreKey(uid uint, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", sessionRefreshTokenStoreKeyPrefix, uid, tokenId)
}

func GenerateSessionRefreshTokenBlocklistStoreKey(tokenId string) string {
	return fmt.Sprintf("%s:%s", sessionRefreshTokenBlocklistKeyPrefix, tokenId)
}

func GetTokenIdFromRefreshSessionKey(tokenKey string) (string, error) {
	if !strings.HasPrefix(tokenKey, sessionRefreshTokenStoreKeyPrefix) {
		return "", errors.New("invalid session key")
	}

	elements := strings.Split(tokenKey, ":")
	if len(elements) != 4 {
		return "", errors.New("invalid session key")
	}

	tokenId := elements[3]
	if tokenId == "" {
		return "", errors.New("invalid session key")
	}
	return tokenId, nil
}

func addSessionToBlocklist(ctx context.Context, sessionKeys []string) {
	for _, tokenKey := range sessionKeys {
		tokenId, err := GetTokenIdFromSessionKey(tokenKey)
		if err != nil {
			log.GetLogger().Warn("Failed to get tokenId from session key", zap.Error(err))
			continue
		}

		// when get the expiresAt failed, use the default expire of the token.
		expire := time.Duration(etc.GetConfig().App.Settings.JWT.TokenExpireTime)
		ttl, err := store.GetRedisClient().GetRDB().TTL(ctx, tokenKey).Result()
		if err == nil {
			// when ttl is less than or equal 5, we consider that this token is about to expire.
			if ttl <= 5 {
				continue
			}
			expire = ttl
		}

		sessionBlocklistKey := GenerateSessionBlocklistKey(tokenId)
		err = store.GetRedisClient().GetRDB().Set(ctx, sessionBlocklistKey, "", expire*time.Second).Err()
		if err != nil {
			log.GetLogger().Warn("Failed to add tokenId to blocklist", zap.Error(err))
		}
	}
}

func addRefreshTokenToBlocklist(ctx context.Context, refreshTokenKeys []string) {
	// add all the refresh tokens of the user to the token blocklist.
	for _, refreshTokenKey := range refreshTokenKeys {
		tokenId, err := GetTokenIdFromRefreshSessionKey(refreshTokenKey)
		if err != nil {
			log.GetLogger().Warn("Failed to get tokenId from session key", zap.Error(err))
			continue
		}

		// when get the expiresAt failed, use the default expire of the token.
		expire := time.Duration(etc.GetConfig().App.Settings.JWT.TokenRefreshTime)
		ttl, err := store.GetRedisClient().GetRDB().TTL(ctx, refreshTokenKey).Result()
		if err == nil {
			// when ttl is less than or equal 5, we consider that this token is about to expire.
			if ttl <= 5 {
				continue
			}
			expire = ttl
		}

		refreshTokenBlocklistKey := GenerateSessionRefreshTokenBlocklistStoreKey(tokenId)
		err = store.GetRedisClient().GetRDB().Set(ctx, refreshTokenBlocklistKey, "", time.Duration(expire)*time.Second).Err()
		if err != nil {
			log.GetLogger().Warn("Failed to add tokenId to blocklist", zap.Error(err))
		}
	}
}

func RevokeUserTokens(ctx context.Context, uid uint) {
	userTokenKeys, err := store.GetRedisClient().ListKeys(ctx, GenerateSessionStoreKey(uid, "*"))
	if err != nil {
		log.GetLogger().Warn("redis scan fail", zap.Error(err))
	}

	addSessionToBlocklist(ctx, userTokenKeys)

	refreshTokenKeys, err := store.GetRedisClient().ListKeys(ctx, GenerateSessionRefreshTokenStoreKey(uid, "*"))
	if err != nil {
		log.GetLogger().Warn("redis scan fail", zap.Error(err))
	}

	addRefreshTokenToBlocklist(ctx, refreshTokenKeys)
}
