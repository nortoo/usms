package encryption

import (
	"fmt"
	"strings"

	"github.com/nortoo/usms/pkg/errors"
	"github.com/nortoo/utils-go/validation"
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

func EncryptEmailAddress(email string) string {
	if !validation.IsValidEmail(email) {
		return ""
	}
	elements := strings.Split(email, ".")
	return fmt.Sprintf("%s***@***.%s", email[:1], elements[len(elements)-1])
}

func EncryptMobileNumber(mobile string) string {
	if isValid, _ := validation.IsValidMobileNumber(mobile, "US"); !isValid {
		return ""
	}

	return fmt.Sprintf("%s****%s", mobile[:3], mobile[len(mobile)-4:])
}
