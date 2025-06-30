package verification

import (
	"context"

	"github.com/nortoo/usms/pkg/errors"
	pb "github.com/nortoo/usms/pkg/proto/verification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedServiceServer
}

func Register(grpcServer *grpc.Server) {
	pb.RegisterServiceServer(grpcServer, &Service{})
}

func (s *Service) ListVerificationMethods(ctx context.Context, req *pb.ListVerificationMethodsReq) (*pb.ListVerificationMethodsResp, error) {
	if req.GetIdentifier() == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}
	return ListVerificationMethods(ctx, req)
}

func (s *Service) GetVerificationTarget(ctx context.Context, req *pb.GetVerificationTargetReq) (*pb.GetVerificationTargetResp, error) {
	if req.GetIdentifier() == "" {
		return nil, errors.ErrInvalidParams.WithDetail("identifier is required")
	}
	if req.GetVerificationMethod() <= 0 {
		return nil, errors.ErrInvalidParams.WithDetail("verification method is required")
	}

	return GetVerificationTarget(ctx, req)
}
