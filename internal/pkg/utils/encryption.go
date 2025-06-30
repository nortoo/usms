package utils

import (
	"github.com/nortoo/usms/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(plainPassword string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.ErrInternalError.WithDetail(err.Error())
	}
	return string(password), nil
}

func ComparePassword(password string, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(plainPassword)) == nil
}
