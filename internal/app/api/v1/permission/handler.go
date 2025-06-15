package permission

import (
	"context"

	pb "github.com/nortoo/usms/pkg/proto/permission/v1"
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

func (s *Service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Permission, error) {
	if req.GetAction() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "action is required.")
	}
	if req.GetResource() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "resource is required.")
	}
	return Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return &emptypb.Empty{}, Delete(ctx, req)
}

func (s *Service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Permission, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return Update(ctx, req)
}

func (s *Service) Get(ctx context.Context, req *pb.GetReq) (*pb.Permission, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return Get(ctx, req)
}

func (s *Service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return List(ctx, req)
}
