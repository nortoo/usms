package group

import (
	"context"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/usergroup/v1"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Group, error)
	Delete(ctx context.Context, req *pb.DeleteReq) error
	Update(ctx context.Context, req *pb.UpdateReq) (*pb.Group, error)
	Get(ctx context.Context, req *pb.GetReq) (*pb.Group, error)
	List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error)
}

type service struct {
	usmCli *usm.Client
}

func NewService(usmCli *usm.Client) Service {
	return &service{usmCli: usmCli}
}

func (s *service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Group, error) {
	g := &model.Group{
		ID:      uint(snowflake.GetSnowWorker().NextId()),
		Name:    req.GetName(),
		Comment: req.GetComment(),
	}
	err := s.usmCli.CreateGroup(g)
	if err != nil {
		return nil, err
	}

	return &pb.Group{
		Id:      uint64(g.ID),
		Name:    g.Name,
		Comment: g.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: g.CreatedAt.Unix(),
			UpdatedAt: g.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteReq) error {
	group := &model.Group{ID: uint(req.GetId())}
	return s.usmCli.DeleteGroup(group)
}

func (s *service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Group, error) {
	g := &model.Group{ID: uint(req.GetId())}

	var cols []string
	if req.GetName() != "" {
		g.Name = req.GetName()
		cols = append(cols, "Name")
	}
	if req.GetComment() != "" {
		g.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}
	if len(cols) == 0 {
		return s.Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := s.usmCli.UpdateGroup(g)
	if err != nil {
		return nil, err
	}

	return &pb.Group{
		Id:      uint64(g.ID),
		Name:    g.Name,
		Comment: g.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: g.CreatedAt.Unix(),
			UpdatedAt: g.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) Get(ctx context.Context, req *pb.GetReq) (*pb.Group, error) {
	g := &model.Group{
		ID:   uint(req.GetId()),
		Name: req.GetName(),
	}
	g, err := s.usmCli.GetGroup(g)
	if err != nil {
		return nil, err
	}

	return &pb.Group{
		Id:      uint64(g.ID),
		Name:    g.Name,
		Comment: g.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: g.CreatedAt.Unix(),
			UpdatedAt: g.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := s.usmCli.ListGroups(&types.QueryGroupOptions{Pagination: &types.Pagination{
		Page:     int(req.GetPagination().GetPage()),
		PageSize: int(req.GetPagination().GetPageSize()),
	}})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Group, len(ret))
	for i, r := range ret {
		items[i] = &pb.Group{
			Id:      uint64(r.ID),
			Name:    r.Name,
			Comment: r.Comment,
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
