package user

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedServiceServer

	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(grpcServer *grpc.Server, handler *Handler) {
	pb.RegisterServiceServer(grpcServer, handler)
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error) {
	if req.GetUsername() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either email or mobile is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return h.service.Create(ctx, req)
}

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return &emptypb.Empty{}, h.service.Delete(ctx, req)
}

func (h *Handler) Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return h.service.Update(ctx, req)
}

func (h *Handler) Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	if req.GetId() <= 0 && req.GetIdentifier() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either uid or identifier is required.")
	}
	return h.service.Get(ctx, req)
}

func (h *Handler) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return h.service.List(ctx, req)
}

func (h *Handler) Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error) {
	if req.GetUsername() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either email or mobile is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return h.service.Signup(ctx, req)
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	if req.GetIdentifier() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return h.service.Login(ctx, req)
}

func (h *Handler) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
	if req.GetToken() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("token is required.")
	}
	if req.GetTenant() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("tenant is required.")
	}
	if req.GetResource() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("resource is required.")
	}
	if req.GetAction() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("action is required.")
	}
	return h.service.Auth(ctx, req)
}

func (h *Handler) DoesIdentifierExist(
	ctx context.Context,
	req *pb.DoesIdentifierExistReq,
) (*pb.DoesIdentifierExistResp, error) {
	if req.GetUsername() == "" && req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either username, email or mobile is required.")
	}
	return h.service.DoesIdentifierExist(ctx, req)
}

func (h *Handler) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*emptypb.Empty, error) {
	if req.GetUid() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("uid is required.")
	}
	if req.GetOldPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("old password is required.")
	}
	if req.GetNewPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("new password is required.")
	}
	return h.service.ChangePassword(ctx, req)
}

func (h *Handler) ResetPassword(ctx context.Context, req *pb.ResetPasswordReq) (*emptypb.Empty, error) {
	if req.GetUid() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("uid is required.")
	}
	if req.GetNewPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("new password is required.")
	}

	return h.service.ResetPassword(ctx, req)
}

func (h *Handler) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenResp, error) {
	if req.GetRefreshToken() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("refresh token is required.")
	}

	return h.service.RefreshToken(ctx, req)
}
