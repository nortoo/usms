package verification

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/verification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *Handler) ListVerificationMethods(
	ctx context.Context,
	req *pb.ListVerificationMethodsReq,
) (*pb.ListVerificationMethodsResp, error) {
	if req.GetIdentifier() == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}
	return h.service.ListVerificationMethods(ctx, req)
}

func (h *Handler) GetVerificationTarget(
	ctx context.Context,
	req *pb.GetVerificationTargetReq,
) (*pb.GetVerificationTargetResp, error) {
	if req.GetIdentifier() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("identifier is required")
	}
	if req.GetVerificationMethod() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("verification method is required")
	}

	return h.service.GetVerificationTarget(ctx, req)
}
