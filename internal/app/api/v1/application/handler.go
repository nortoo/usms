package application

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/application/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pb.UnimplementedServiceServer
}

func Register(grpcServer *grpc.Server) {
	pb.RegisterServiceServer(grpcServer, &Service{})
}

func (s *Service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Application, error) {
	if req.GetName() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("name is required.")
	}
	return Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return &emptypb.Empty{}, Delete(ctx, req)
}

func (s *Service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Application, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return Update(ctx, req)
}

func (s *Service) Get(ctx context.Context, req *pb.GetReq) (*pb.Application, error) {
	if req.GetId() <= 0 &&
		req.GetName() == "" &&
		req.GetAppid() == "" {
		return nil,
			errors.ErrInvalidParams.WithDetail("query an application required at least one condition.")
	}
	return Get(ctx, req)
}

func (s *Service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return List(ctx, req)
}
