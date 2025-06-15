package group

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/usergroup/v1"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.Group, error) {
	g := &model.Group{
		ID:      uint(snowflake.GetSnowWorker().NextId()),
		Name:    req.GetName(),
		Comment: req.GetComment(),
	}
	err := _usm.Client().CreateGroup(g)
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

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	group := &model.Group{ID: uint(req.GetId())}
	return _usm.Client().DeleteGroup(group)
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.Group, error) {
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
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := _usm.Client().UpdateGroup(g)
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

func Get(ctx context.Context, req *pb.GetReq) (*pb.Group, error) {
	g := &model.Group{
		ID:   uint(req.GetId()),
		Name: req.GetName(),
	}
	g, err := _usm.Client().GetGroup(g)
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

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := _usm.Client().ListGroups(&types.QueryGroupOptions{Pagination: &types.Pagination{
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
