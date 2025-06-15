package user

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/internal/pkg/types/user"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"gorm.io/gorm"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error) {
	u := &model.User{
		Model:    gorm.Model{ID: uint(snowflake.GetSnowWorker().NextId())},
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		State:    int8(req.GetState()),
	}

	for _, rid := range req.GetRoles() {
		role, err := _usm.Client().GetRole(&model.Role{ID: uint(rid)})
		if err != nil {
			continue
		}
		u.Roles = append(u.Roles, role)
	}

	for _, gid := range req.GetGroups() {
		group, err := _usm.Client().GetGroup(&model.Group{ID: uint(gid)})
		if err != nil {
			continue
		}
		u.Groups = append(u.Groups, group)
	}

	err := _usm.Client().CreateUser(u)
	if err != nil {
		return nil, err
	}

	return user.ModelToPb(u), nil
}

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	return _usm.Client().DeleteUser(&model.User{Model: gorm.Model{ID: uint(req.GetId())}})
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	u := &model.User{Model: gorm.Model{ID: uint(req.GetId())}}

	var cols []string
	if req.GetMobile() != "" {
		u.Mobile = req.GetMobile()
		cols = append(cols, "Mobile")
	}
	if req.GetEmail() != "" {
		u.Email = req.GetEmail()
		cols = append(cols, "Email")
	}
	if req.GetPassword() != "" {
		u.Password = req.GetPassword()
	}
	if req.GetState() != 0 {
		u.State = int8(req.GetState())
	}

	var roles []*model.Role
	for _, rid := range req.GetRoles() {
		role, err := _usm.Client().GetRole(&model.Role{ID: uint(rid)})
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}
	if len(roles) != 0 {
		u.Roles = roles
		cols = append(cols, "Roles")
	}

	var groups []*model.Group
	for _, gid := range req.GetGroups() {
		group, err := _usm.Client().GetGroup(&model.Group{ID: uint(gid)})
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}
	if len(groups) != 0 {
		u.Groups = groups
		cols = append(cols, "Groups")
	}

	if len(cols) == 0 {
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	return user.ModelToPb(u), nil
}

func Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	u, err := _usm.Client().GetUser(&model.User{Model: gorm.Model{ID: uint(req.GetId())}})
	if err != nil {
		return nil, err
	}
	return user.ModelToPb(u), nil
}

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	options := &types.QueryUserOptions{
		Pagination: &types.Pagination{
			Page:     int(req.GetPagination().GetPage()),
			PageSize: int(req.GetPagination().GetPageSize()),
		},
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		GroupID:  uint(req.GetGroupId()),
		RoleID:   uint(req.GetRoleId()),
	}
	for _, state := range req.GetState() {
		options.States = append(options.States, int8(state))
	}
	ret, total, err := _usm.Client().ListUsers(options)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.User, len(ret))
	for i, r := range ret {
		items[i] = user.ModelToPb(r)
	}

	return &pb.ListResp{
		Pagination: &pbtypes.PaginationResp{
			Page:  req.GetPagination().Page,
			Total: total,
		},
		Items: items,
	}, nil
}
