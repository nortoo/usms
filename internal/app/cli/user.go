package cli

import (
	"context"
	"fmt"
	"time"

	applicationpb "github.com/nortoo/usms/pkg/proto/application/v1"
	permissionpb "github.com/nortoo/usms/pkg/proto/permission/v1"
	rolepb "github.com/nortoo/usms/pkg/proto/role/v1"
	userpb "github.com/nortoo/usms/pkg/proto/user/v1"
	"github.com/pkg/errors"
)

func (c *Client) CreateAdministration(username, password, tenant, email string, createAll bool) (*userpb.User, error) {
	userCli := userpb.NewServiceClient(c.conn)
	roleCli := rolepb.NewServiceClient(c.conn)
	applicationCli := applicationpb.NewServiceClient(c.conn)
	permissionCli := permissionpb.NewServiceClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app, err := applicationCli.Get(ctx, &applicationpb.GetReq{Name: tenant})
	if err != nil {
		if createAll {
			app, err = applicationCli.Create(ctx, &applicationpb.CreateReq{Name: tenant})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("tenant [%s] does not exist, try to use '-create-all true' to create all not existing entities.", tenant)
		}
	}

	permission, err := permissionCli.Get(ctx, &permissionpb.GetReq{
		Action:   "*",
		Resource: "*",
	})
	if err != nil {
		fmt.Println(err)
		if createAll {
			permission, err = permissionCli.Create(ctx, &permissionpb.CreateReq{
				Action:   "*",
				Resource: "*",
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("default super permission does not exist, try to use '-create-all true' to create all not existing entities.")
		}
	}

	role, err := roleCli.Get(ctx, &rolepb.GetReq{Name: "superadmin"})
	if err != nil {
		if createAll {
			role, err = roleCli.Create(ctx, &rolepb.CreateReq{
				Name:          "superadmin",
				Comment:       "this is a super administrator.",
				ApplicationId: app.GetId(),
				Permissions:   []uint64{permission.GetId()},
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("tenant [%s] does not exist, try to use '-create-all true' to create all not existing entities.", tenant)
		}
	}

	return userCli.Create(ctx, &userpb.CreateReq{
		Username: username,
		Password: password,
		Email:    email,
		Roles:    []int64{int64(role.Id)},
		State:    1,
	})
}
