package permission

import (
	"context"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/permission/v1"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Permission, error)
	Delete(ctx context.Context, req *pb.DeleteReq) error
	Update(ctx context.Context, req *pb.UpdateReq) (*pb.Permission, error)
	Get(ctx context.Context, req *pb.GetReq) (*pb.Permission, error)
	List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error)
}

type service struct {
	usmCli *usm.Client
}

func NewService(usmCli *usm.Client) Service {
	return &service{usmCli: usmCli}
}

func (s *service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Permission, error) {
	p := &model.Permission{
		ID:       uint(snowflake.GetSnowWorker().NextId()),
		Action:   req.GetAction(),
		Resource: req.GetResource(),
		Comment:  req.GetComment(),
	}
	err := s.usmCli.CreatePermission(p)
	if err != nil {
		return nil, err
	}

	return &pb.Permission{
		Id:       uint64(p.ID),
		Action:   p.Action,
		Resource: p.Resource,
		Comment:  p.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteReq) error {
	return s.usmCli.DeletePermission(&model.Permission{ID: uint(req.GetId())})
}

func (s *service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Permission, error) {
	p := &model.Permission{ID: uint(req.GetId())}

	var cols []string
	if req.GetComment() != "" {
		p.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}
	if len(cols) == 0 {
		return s.Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := s.usmCli.UpdatePermission(p)
	if err != nil {
		return nil, err
	}

	return &pb.Permission{
		Id:       uint64(p.ID),
		Action:   p.Action,
		Resource: p.Resource,
		Comment:  p.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Get(ctx context.Context, req *pb.GetReq) (*pb.Permission, error) {
	p, err := s.usmCli.GetPermission(&model.Permission{
		ID:       uint(req.GetId()),
		Resource: req.GetResource(),
		Action:   req.GetAction(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.Permission{
		Id:       uint64(p.ID),
		Action:   p.Action,
		Resource: p.Resource,
		Comment:  p.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := s.usmCli.ListPermissions(&types.QueryPermissionOptions{Pagination: &types.Pagination{
		Page:     int(req.GetPagination().GetPage()),
		PageSize: int(req.GetPagination().GetPageSize()),
	}})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Permission, len(ret))
	for i, r := range ret {
		items[i] = &pb.Permission{
			Id:       uint64(r.ID),
			Action:   r.Action,
			Resource: r.Resource,
			Comment:  r.Comment,
			Time: &pbtypes.TimeModel{
				CreatedAt: r.CreatedAt.Unix(),
				UpdatedAt: r.UpdatedAt.Unix(),
			},
		}
	}

	return &pb.ListResp{
		Pagination: &pbtypes.PaginationResp{
			Page:  req.GetPagination().Page,
			Total: total,
		},
		Items: items,
	}, nil
}
