package permission

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/permission/v1"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.Permission, error) {
	p := &model.Permission{
		ID:       uint(snowflake.GetSnowWorker().NextId()),
		Action:   req.GetAction(),
		Resource: req.GetResource(),
		Comment:  req.GetComment(),
	}
	err := _usm.Client().CreatePermission(p)
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

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	return _usm.Client().DeletePermission(&model.Permission{ID: uint(req.GetId())})
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.Permission, error) {
	p := &model.Permission{ID: uint(req.GetId())}

	var cols []string
	if req.GetComment() != "" {
		p.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}
	if len(cols) == 0 {
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := _usm.Client().UpdatePermission(p)
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

func Get(ctx context.Context, req *pb.GetReq) (*pb.Permission, error) {
	p, err := _usm.Client().GetPermission(&model.Permission{ID: uint(req.GetId())})
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

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := _usm.Client().ListPermissions(&types.QueryPermissionOptions{Pagination: &types.Pagination{
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
