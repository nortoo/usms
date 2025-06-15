package errors

import (
	"google.golang.org/grpc/codes"
)

type USMSError struct {
	Code    codes.Code
	Message string
}

func newUSMSError(code codes.Code, message string) *USMSError {
	return &USMSError{Code: code, Message: message}
}

func (e *USMSError) WithDetail(message string) *USMSError {
	return &USMSError{
		Code:    e.Code,
		Message: message,
	}
}

func (e *USMSError) Error() string {
	return e.Message
}

var (
	ErrInvalidParams    = newUSMSError(codes.InvalidArgument, "invalid params")
	ErrLoginDeclined    = newUSMSError(codes.Unauthenticated, "username or password is incorrect")
	ErrUnauthenticated  = newUSMSError(codes.Unauthenticated, "user is not authenticated")
	ErrPermissionDenied = newUSMSError(codes.PermissionDenied, "permission denied")

	ErrResourceNotFound = newUSMSError(codes.NotFound, "resource does not exist")
)
