package menu

import (
	"context"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/pkg/errors"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/menu/v1"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Menu, error)
	Delete(ctx context.Context, req *pb.DeleteReq) error
	Update(ctx context.Context, req *pb.UpdateReq) (*pb.Menu, error)
	Get(ctx context.Context, req *pb.GetReq) (*pb.Menu, error)
	List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error)
}

type service struct {
	usmCli *usm.Client
}

func NewService(usmCli *usm.Client) Service {
	return &service{usmCli: usmCli}
}

func (s *service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Menu, error) {
	if req.GetParentId() != 0 {
		_, err := s.usmCli.GetMenu(&model.Menu{ID: uint(req.GetParentId())})
		if err != nil {
			return nil, errors.ErrResourceNotFound.WithDetail("parent menu does not exist")
		}
	}
	m := &model.Menu{
		ID:       uint(snowflake.GetSnowWorker().NextId()),
		ParentID: int64(req.GetParentId()),
		Name:     req.GetName(),
		Path:     req.GetPath(),
		Comment:  req.GetComment(),
	}
	err := s.usmCli.CreateMenu(m)
	if err != nil {
		return nil, err
	}

	return &pb.Menu{
		Id:       uint64(m.ID),
		ParentId: uint64(m.ParentID),
		Name:     m.Name,
		Path:     m.Path,
		Comment:  m.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: m.CreatedAt.Unix(),
			UpdatedAt: m.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteReq) error {
	return s.usmCli.DeleteMenu(&model.Menu{ID: uint(req.GetId())})
}

func (s *service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Menu, error) {
	m := &model.Menu{ID: uint(req.GetId())}

	var cols []string
	if req.GetName() != "" {
		m.Name = req.GetName()
		cols = append(cols, "Name")
	}
	if req.GetPath() != "" {
		m.Path = req.GetPath()
		cols = append(cols, "Path")
	}
	if req.GetComment() != "" {
		m.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}
	if len(cols) == 0 {
		return s.Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := s.usmCli.UpdateMenu(m)
	if err != nil {
		return nil, err
	}

	return &pb.Menu{
		Id:       uint64(m.ID),
		ParentId: uint64(m.ParentID),
		Name:     m.Name,
		Path:     m.Path,
		Comment:  m.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: m.CreatedAt.Unix(),
			UpdatedAt: m.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Get(ctx context.Context, req *pb.GetReq) (*pb.Menu, error) {
	m := &model.Menu{ID: uint(req.GetId())}
	m, err := s.usmCli.GetMenu(m)
	if err != nil {
		return nil, err
	}

	return &pb.Menu{
		Id:       uint64(m.ID),
		ParentId: uint64(m.ParentID),
		Name:     m.Name,
		Path:     m.Path,
		Comment:  m.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: m.CreatedAt.Unix(),
			UpdatedAt: m.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := s.usmCli.ListMenus(&types.QueryMenuOptions{Pagination: &types.Pagination{
		Page:     int(req.GetPagination().GetPage()),
		PageSize: int(req.GetPagination().GetPageSize()),
	}})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Menu, len(ret))
	for i, r := range ret {
		items[i] = &pb.Menu{
			Id:       uint64(r.ID),
			ParentId: uint64(r.ParentID),
			Name:     r.Name,
			Path:     r.Path,
			Comment:  r.Comment,
			Time: &pbtypes.TimeModel{
				CreatedAt: r.CreatedAt.Unix(),
				UpdatedAt: r.UpdatedAt.Unix(),
			},
		}
	}

	return &pb.ListResp{
		Pagination: &pbtypes.PaginationResp{
			Page:  req.GetPagination().GetPage(),
			Total: total,
		},
		Items: items,
	}, nil
}
