package validation

import (
	"fmt"
	"regexp"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/pkg/errors"
)

func IsValidUsername(username string) (bool, error) {
	minLength := etc.GetConfig().App.Settings.Validation.UsernamePattern.MinLength
	if minLength <= 0 {
		minLength = 5
	}
	maxLength := etc.GetConfig().App.Settings.Validation.UsernamePattern.MaxLength
	if maxLength <= 0 {
		maxLength = 20
	}

	pattern := regexp.MustCompile(fmt.Sprintf("^[a-zA-Z][a-zA-Z0-9_]{%d,%d}$", minLength-1, maxLength-1))
	if !pattern.MatchString(username) {
		return false, errors.Errorf("username only accepts %d-%d characters, alphanumeric and underscores, and should start with a letter.", minLength, maxLength)
	}
	return true, nil
}

func IsValidPassword(password string) (bool, error) {
	minLength := etc.GetConfig().App.Settings.Validation.PasswordPattern.MinLength
	if minLength <= 0 {
		minLength = 8
	}

	lengthRegex := regexp.MustCompile(fmt.Sprintf(`.{%d,}`, minLength))
	if !lengthRegex.MatchString(password) {
		return false, errors.Errorf("password must be at least %d characters", minLength)
	}

	if etc.GetConfig().App.Settings.Validation.PasswordPattern.RequireLowerCase {
		lower := regexp.MustCompile(`[a-z]`)
		if !lower.MatchString(password) {
			return false, errors.New("password must contain at least one lowercase character")
		}
	}

	if etc.GetConfig().App.Settings.Validation.PasswordPattern.RequireUpperCase {
		upper := regexp.MustCompile(`[A-Z]`)
		if !upper.MatchString(password) {
			return false, errors.New("password must contain at least one uppercase character")
		}
	}

	if etc.GetConfig().App.Settings.Validation.PasswordPattern.RequireDigit {
		digit := regexp.MustCompile(`[0-9]`)
		if !digit.MatchString(password) {
			return false, errors.New("password must contain at least one digit character")
		}
	}

	if etc.GetConfig().App.Settings.Validation.PasswordPattern.RequireSpecialChars {
		special := regexp.MustCompile(`[` + "`" + `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`)
		if !special.MatchString(password) {
			return false, errors.New("password must contain at least one special character")
		}
	}

	return true, nil
}
