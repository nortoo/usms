package role

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/internal/pkg/types/role"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/role/v1"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.Role, error) {
	r := &model.Role{
		ID:            uint(snowflake.GetSnowWorker().NextId()),
		Name:          req.GetName(),
		Comment:       req.GetComment(),
		ApplicationID: uint(req.GetApplicationId()),
	}
	for _, mid := range req.GetMenus() {
		m, err := _usm.Client().GetMenu(&model.Menu{ID: uint(mid)})
		if err != nil {
			continue
		}
		r.Menus = append(r.Menus, m)
	}
	for _, pid := range req.GetPermissions() {
		p, err := _usm.Client().GetPermission(&model.Permission{ID: uint(pid)})
		if err != nil {
			continue
		}
		r.Permissions = append(r.Permissions, p)
	}

	err := _usm.Client().CreateRole(r)
	if err != nil {
		return nil, err
	}

	return role.ModelToPb(r), nil
}

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	return _usm.Client().DeleteRole(&model.Role{ID: uint(req.GetId())})
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.Role, error) {
	r := &model.Role{ID: uint(req.GetId())}

	var cols []string
	if req.GetComment() != "" {
		r.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}

	var menus []*model.Menu
	for _, mid := range req.GetMenus() {
		m, err := _usm.Client().GetMenu(&model.Menu{ID: uint(mid)})
		if err != nil {
			continue
		}
		menus = append(menus, m)
	}
	if len(menus) != 0 {
		r.Menus = menus
		cols = append(cols, "Menus")
	}

	var permissions []*model.Permission
	for _, pid := range req.GetPermissions() {
		p, err := _usm.Client().GetPermission(&model.Permission{ID: uint(pid)})
		if err != nil {
			continue
		}
		permissions = append(permissions, p)
	}
	if len(permissions) != 0 {
		r.Permissions = permissions
		cols = append(cols, "Permissions")
	}

	if len(cols) == 0 {
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := _usm.Client().UpdateRole(r)
	if err != nil {
		return nil, err
	}

	return role.ModelToPb(r), nil
}

func Get(ctx context.Context, req *pb.GetReq) (*pb.Role, error) {
	r, err := _usm.Client().GetRole(&model.Role{ID: uint(req.GetId())})
	if err != nil {
		return nil, err
	}

	return role.ModelToPb(r), nil
}

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := _usm.Client().ListRoles(&types.QueryRoleOptions{Pagination: &types.Pagination{
		Page:     int(req.GetPagination().GetPage()),
		PageSize: int(req.GetPagination().GetPageSize()),
	}})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Role, len(ret))
	for i, r := range ret {
		items[i] = &pb.Role{
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
			Page:  req.GetPagination().Page,
			Total: total,
		},
		Items: items,
	}, nil
}
