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
