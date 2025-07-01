package verification

import (
	"context"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usms/internal/pkg/utils/encryption"
	"github.com/nortoo/usms/internal/pkg/utils/identification"
	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/verification/v1"
)

type Service interface {
	ListVerificationMethods(ctx context.Context, req *pb.ListVerificationMethodsReq) (*pb.ListVerificationMethodsResp, error)
	GetVerificationTarget(ctx context.Context, req *pb.GetVerificationTargetReq) (*pb.GetVerificationTargetResp, error)
}

type service struct {
	usmCli            *usm.Client
	identificationSvc identification.Service
}

func NewService(usmCli *usm.Client, identificationSvc identification.Service) Service {
	return &service{
		usmCli:            usmCli,
		identificationSvc: identificationSvc,
	}
}

// ListVerificationMethods lists the available verification methods through a user's identifier.
// The `identifier` could be either username, email, or mobile.
func (s *service) ListVerificationMethods(ctx context.Context, req *pb.ListVerificationMethodsReq) (*pb.ListVerificationMethodsResp, error) {
	var u *model.User
	var err error

	identifier := req.GetIdentifier()
	switch s.identificationSvc.Recognize(identifier) {
	case identification.Email:
		u, err = s.usmCli.GetUser(&model.User{Email: identifier})
	case identification.Mobile:
		u, err = s.usmCli.GetUser(&model.User{Mobile: identifier})
	case identification.Username:
		u, err = s.usmCli.GetUser(&model.User{Username: identifier})
	default:
		return nil, errors.ErrInvalidParams.WithDetail("invalid identifier")
	}
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	availableVerificationMethods := make([]*pb.VerificationMethod, 0)
	if u.Email != "" {
		availableVerificationMethods = append(availableVerificationMethods, &pb.VerificationMethod{
			VerificationMethod: pb.EnumVerification_Email,
			Target:             encryption.EncryptEmailAddress(u.Email),
		})
	}
	if u.Mobile != "" {
		availableVerificationMethods = append(availableVerificationMethods, &pb.VerificationMethod{
			VerificationMethod: pb.EnumVerification_Mobile,
			Target:             encryption.EncryptMobileNumber(u.Mobile),
		})
	}

	return &pb.ListVerificationMethodsResp{AvailableVerificationMethods: availableVerificationMethods}, nil
}

func (s *service) GetVerificationTarget(ctx context.Context, req *pb.GetVerificationTargetReq) (*pb.GetVerificationTargetResp, error) {
	var u *model.User
	var err error

	identifier := req.GetIdentifier()
	switch s.identificationSvc.Recognize(identifier) {
	case identification.Email:
		u, err = s.usmCli.GetUser(&model.User{Email: identifier})
	case identification.Mobile:
		u, err = s.usmCli.GetUser(&model.User{Mobile: identifier})
	case identification.Username:
		u, err = s.usmCli.GetUser(&model.User{Username: identifier})
	default:
		return nil, errors.ErrInvalidParams.WithDetail("invalid identifier")
	}
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	var target string
	switch req.GetVerificationMethod() {
	case pb.EnumVerification_Email:
		target = u.Email
	case pb.EnumVerification_Mobile:
		target = u.Mobile
	default:
		return nil, errors.ErrInvalidParams.WithDetail("invalid verification method")
	}

	if target == "" {
		return nil, errors.ErrResourceNotFound.WithDetail("verification target not found")
	}

	return &pb.GetVerificationTargetResp{
		VerificationMethod: req.GetVerificationMethod(),
		Target:             target,
	}, nil
}
