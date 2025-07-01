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

type Service interface {
	Recognize(identifier string) Identifier
}

type service struct {
	validator _validation.Service
}

func New(validator _validation.Service) Service {
	return &service{
		validator: validator,
	}
}

func (s *service) Recognize(identifier string) Identifier {
	if validation.IsValidEmail(identifier) {
		return Email
	} else if isValid, _ := validation.IsValidMobileNumber(identifier, "US"); isValid {
		return Mobile
	} else if _, err := s.validator.IsValidUsername(identifier); err == nil {
		return Username
	} else {
		return "unknown"
	}
}
