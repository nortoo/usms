package role

import (
	"context"

	pb "github.com/nortoo/usms/pkg/proto/role/v1"
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

func (s *Service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Role, error) {
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required.")
	}
	if req.GetApplicationId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "application_id is required.")
	}
	return Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return &emptypb.Empty{}, Delete(ctx, req)
}

func (s *Service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Role, error) {
	if req.GetId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required.")
	}
	return Update(ctx, req)
}

func (s *Service) Get(ctx context.Context, req *pb.GetReq) (*pb.Role, error) {
	if req.GetId() <= 0 &&
		req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "query a role required at least one condition.")
	}
	return Get(ctx, req)
}

func (s *Service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return List(ctx, req)
}
