package role

import (
	"github.com/nortoo/usm/model"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	menupb "github.com/nortoo/usms/pkg/proto/menu/v1"
	permissionpb "github.com/nortoo/usms/pkg/proto/permission/v1"
	pb "github.com/nortoo/usms/pkg/proto/role/v1"
)

func ModelToPb(m *model.Role) *pb.Role {
	p := &pb.Role{
		Id:      uint64(m.ID),
		Name:    m.Name,
		Comment: m.Comment,
		Time: &pbtypes.TimeModel{
			CreatedAt: m.CreatedAt.Unix(),
			UpdatedAt: m.UpdatedAt.Unix(),
		},
	}

	if m.Application != nil {
		p.Application = m.Application.Name
	}

	for _, menu := range m.Menus {
		p.Menus = append(p.Menus, &menupb.Menu{
			Id:       uint64(menu.ID),
			ParentId: uint64(menu.ParentID),
			Name:     menu.Name,
			Path:     menu.Path,
			Comment:  menu.Comment,
			Time: &pbtypes.TimeModel{
				CreatedAt: menu.CreatedAt.Unix(),
				UpdatedAt: menu.UpdatedAt.Unix(),
			},
		})
	}

	for _, permission := range m.Permissions {
		p.Permissions = append(p.Permissions, &permissionpb.Permission{
			Id:       uint64(permission.ID),
			Action:   permission.Action,
			Resource: permission.Resource,
			Comment:  permission.Comment,
			Time: &pbtypes.TimeModel{
				CreatedAt: permission.CreatedAt.Unix(),
				UpdatedAt: permission.UpdatedAt.Unix(),
			},
		})
	}

	return p
}
