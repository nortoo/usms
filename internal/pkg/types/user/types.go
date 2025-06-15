package user

import (
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usms/internal/pkg/types/role"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	grouppb "github.com/nortoo/usms/pkg/proto/usergroup/v1"
)

func ModelToPb(m *model.User) *pb.User {
	p := &pb.User{
		Id:       uint64(m.ID),
		Username: m.Username,
		Email:    m.Email,
		Mobile:   m.Mobile,
		State:    int32(m.State),
		Time: &pbtypes.TimeModel{
			CreatedAt: m.CreatedAt.Unix(),
			UpdatedAt: m.UpdatedAt.Unix(),
		},
	}

	for _, r := range m.Roles {
		p.Roles = append(p.Roles, role.ModelToPb(r))
	}
	for _, g := range m.Groups {
		p.Groups = append(p.Groups, &grouppb.Group{
			Id:      uint64(g.ID),
			Name:    g.Name,
			Comment: g.Comment,
			Time: &pbtypes.TimeModel{
				CreatedAt: g.CreatedAt.Unix(),
				UpdatedAt: g.UpdatedAt.Unix(),
			},
		})
	}

	return p
}
