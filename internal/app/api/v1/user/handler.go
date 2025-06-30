package user

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedServiceServer
}

func Register(grpcServer *grpc.Server) {
	pb.RegisterServiceServer(grpcServer, &Service{})
}

func (s *Service) Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error) {
	if req.GetUsername() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either email or mobile is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return &emptypb.Empty{}, Delete(ctx, req)
}

func (s *Service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return Update(ctx, req)
}

func (s *Service) Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return Get(ctx, req)
}

func (s *Service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return List(ctx, req)
}

func (s *Service) Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error) {
	if req.GetUsername() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either email or mobile is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return Signup(ctx, req)
}

func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	if req.GetIdentifier() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("username is required.")
	}
	if req.GetPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("password is required.")
	}
	return Login(ctx, req)
}

func (s *Service) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
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
	return Auth(ctx, req)
}

func (s *Service) DoesIdentifierExist(ctx context.Context, req *pb.DoesIdentifierExistReq) (*pb.DoesIdentifierExistResp, error) {
	if req.GetUsername() == "" && req.GetEmail() == "" && req.GetMobile() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("either username, email or mobile is required.")
	}
	return DoesIdentifierExist(ctx, req)
}

func (s *Service) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*emptypb.Empty, error) {
	if req.GetUid() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("uid is required.")
	}
	if req.GetOldPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("old password is required.")
	}
	if req.GetNewPassword() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("new password is required.")
	}
	return ChangePassword(ctx, req)
}
