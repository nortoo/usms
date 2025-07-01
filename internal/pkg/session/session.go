package session

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service interface {
	GenerateSessionStoreKey(uid uint, tokenId string) string
	GenerateSessionBlocklistKey(tokenId string) string
	GetTokenIdFromSessionKey(tokenKey string) (string, error)
	GenerateSessionRefreshTokenStoreKey(uid uint, tokenId string) string
	GenerateSessionRefreshTokenBlocklistStoreKey(tokenId string) string
	GetTokenIdFromRefreshSessionKey(tokenKey string) (string, error)
	addSessionToBlocklist(ctx context.Context, sessionKeys []string)
	addRefreshTokenToBlocklist(ctx context.Context, refreshTokenKeys []string)
	RevokeUserTokens(ctx context.Context, uid uint)
}

type service struct {
	config   *etc.Config
	redisCli *store.RedisCli
	logger   *zap.Logger
}

func NewService(conf *etc.Config, redisCli *store.RedisCli, logger *zap.Logger) Service {
	return &service{
		config:   conf,
		redisCli: redisCli,
		logger:   logger,
	}
}

const (
	sessionStoreKeyPrefix             = "session:user"
	sessionRefreshTokenStoreKeyPrefix = "session:refresh"

	sessionBlocklistKeyPrefix             = "session:blocklist:token"
	sessionRefreshTokenBlocklistKeyPrefix = "session:blocklist:refresh"
)

func (s *service) GenerateSessionStoreKey(uid uint, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", sessionStoreKeyPrefix, uid, tokenId)
}

func (s *service) GenerateSessionBlocklistKey(tokenId string) string {
	return fmt.Sprintf("%s:%s", sessionBlocklistKeyPrefix, tokenId)
}

func (s *service) GetTokenIdFromSessionKey(tokenKey string) (string, error) {
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

func (s *service) GenerateSessionRefreshTokenStoreKey(uid uint, tokenId string) string {
	return fmt.Sprintf("%s:%d:%s", sessionRefreshTokenStoreKeyPrefix, uid, tokenId)
}

func (s *service) GenerateSessionRefreshTokenBlocklistStoreKey(tokenId string) string {
	return fmt.Sprintf("%s:%s", sessionRefreshTokenBlocklistKeyPrefix, tokenId)
}

func (s *service) GetTokenIdFromRefreshSessionKey(tokenKey string) (string, error) {
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

func (s *service) addSessionToBlocklist(ctx context.Context, sessionKeys []string) {
	for _, tokenKey := range sessionKeys {
		tokenId, err := s.GetTokenIdFromSessionKey(tokenKey)
		if err != nil {
			s.logger.Warn("Failed to get tokenId from session key", zap.Error(err))
			continue
		}

		// when get the expiresAt failed, use the default expire of the token.
		expire := time.Duration(s.config.App.Settings.JWT.TokenExpireTime) * time.Second
		ttl, err := s.redisCli.GetRDB().TTL(ctx, tokenKey).Result()
		if err == nil {
			// when ttl is less than or equal 5, we consider that this token is about to expire.
			if ttl <= 5 {
				continue
			}
			expire = ttl
		}

		sessionBlocklistKey := s.GenerateSessionBlocklistKey(tokenId)
		err = s.redisCli.GetRDB().Set(ctx, sessionBlocklistKey, "", expire).Err()
		if err != nil {
			s.logger.Warn("Failed to add tokenId to blocklist", zap.Error(err))
		}
	}
}

func (s *service) addRefreshTokenToBlocklist(ctx context.Context, refreshTokenKeys []string) {
	// add all the refresh tokens of the user to the token blocklist.
	for _, refreshTokenKey := range refreshTokenKeys {
		tokenId, err := s.GetTokenIdFromRefreshSessionKey(refreshTokenKey)
		if err != nil {
			s.logger.Warn("Failed to get tokenId from session key", zap.Error(err))
			continue
		}

		// when get the expiresAt failed, use the default expire of the token.
		expire := time.Duration(s.config.App.Settings.JWT.TokenRefreshTime) * time.Second
		ttl, err := s.redisCli.GetRDB().TTL(ctx, refreshTokenKey).Result()
		if err == nil {
			// when ttl is less than or equal 5, we consider that this token is about to expire.
			if ttl <= 5 {
				continue
			}
			expire = ttl
		}

		refreshTokenBlocklistKey := s.GenerateSessionRefreshTokenBlocklistStoreKey(tokenId)
		err = s.redisCli.GetRDB().Set(ctx, refreshTokenBlocklistKey, "", expire).Err()
		if err != nil {
			s.logger.Warn("Failed to add tokenId to blocklist", zap.Error(err))
		}
	}
}

func (s *service) RevokeUserTokens(ctx context.Context, uid uint) {
	userTokenKeys, err := s.redisCli.ListKeys(ctx, s.GenerateSessionStoreKey(uid, "*"))
	if err != nil {
		s.logger.Warn("redis scan fail", zap.Error(err))
	}

	s.addSessionToBlocklist(ctx, userTokenKeys)

	refreshTokenKeys, err := s.redisCli.ListKeys(ctx, s.GenerateSessionRefreshTokenStoreKey(uid, "*"))
	if err != nil {
		s.logger.Warn("redis scan fail", zap.Error(err))
	}

	s.addRefreshTokenToBlocklist(ctx, refreshTokenKeys)
}
