package session

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func GenerateSessionStoreKey(uid uint, tokenId string) string {
	return fmt.Sprintf("session:user:%d:%s", uid, tokenId)
}

func GenerateSessionBlocklistKey(tokenId string) string {
	return fmt.Sprintf("session:blocklist:%s", tokenId)
}

func GetTokenIdFromSessionKey(tokenKey string) (string, error) {
	if !strings.HasPrefix(tokenKey, "session:user:") {
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
