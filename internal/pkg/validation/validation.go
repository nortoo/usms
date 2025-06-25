package validation

import (
	"regexp"

	"github.com/nortoo/usms/internal/pkg/etc"
)

func IsValidUsername(username string) bool {
	pattern := etc.GetConfig().App.Settings.UsernamePattern
	reg, _ := regexp.Compile(pattern)
	return reg.MatchString(username)
}

func IsValidPassword(password string) bool {
	pattern := etc.GetConfig().App.Settings.PasswordPattern
	reg, _ := regexp.Compile(pattern)
	return reg.MatchString(password)
}
