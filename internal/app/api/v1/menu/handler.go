package menu

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/menu/v1"
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

func (h *Handler) Create(ctx context.Context, req *pb.CreateReq) (*pb.Menu, error) {
	if req.GetName() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("name is required.")
	}
	if req.GetPath() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("path is required.")
	}
	return h.service.Create(ctx, req)
}

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteReq) (*emptypb.Empty, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return &emptypb.Empty{}, h.service.Delete(ctx, req)
}

func (h *Handler) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Menu, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return h.service.Update(ctx, req)
}

func (h *Handler) Get(ctx context.Context, req *pb.GetReq) (*pb.Menu, error) {
	if req.GetId() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("id is required.")
	}
	return h.service.Get(ctx, req)
}

func (h *Handler) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	return h.service.List(ctx, req)
}
