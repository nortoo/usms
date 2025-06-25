package user

import (
	"context"

	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.InvalidArgument, "username is required.")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required.")
	}
	if req.GetState() == 0 {
		return nil, status.Error(codes.InvalidArgument, "state is required.")
	}
	return Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return &emptypb.Empty{}, Delete(ctx, req)
}

func (s *Service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return Update(ctx, req)
}

func (s *Service) Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return Get(ctx, req)
}

func (s *Service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return List(ctx, req)
}

func (s *Service) Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error) {
	if req.GetUsername() == "" ||
		req.GetEmail() == "" ||
		req.GetMobile() == "" {
		return nil, status.Error(codes.InvalidArgument, "either username, email or mobile is required.")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required.")
	}
	return Signup(ctx, req)
}

func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	if req.GetIdentifier() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required.")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required.")
	}
	return Login(ctx, req)
}

func (s *Service) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required.")
	}
	if req.GetTenant() == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant is required.")
	}
	if req.GetResource() == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required.")
	}
	if req.GetAction() == "" {
		return nil, status.Error(codes.InvalidArgument, "action is required.")
	}
	return Auth(ctx, req)
}
