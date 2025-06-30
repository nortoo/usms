package identification

import (
	_validation "github.com/nortoo/usms/internal/pkg/validation"
	"github.com/nortoo/utils-go/validation"
)

// The Identifier represents the identity type for login including username, email, and mobile.
type Identifier string

const (
	Username Identifier = "username"
	Email    Identifier = "email"
	Mobile   Identifier = "mobile"
)

func Recognize(identifier string) Identifier {
	if validation.IsValidEmail(identifier) {
		return Email
	} else if isValid, _ := validation.IsValidMobileNumber(identifier, "US"); isValid {
		return Mobile
	} else if _, err := _validation.IsValidUsername(identifier); err == nil {
		return Username
	} else {
		return "unknown"
	}
}
