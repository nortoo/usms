package test

import (
	"context"
	"testing"
	"time"

	applicationpb "github.com/nortoo/usms/pkg/proto/application/v1"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	menupb "github.com/nortoo/usms/pkg/proto/menu/v1"
	permissionpb "github.com/nortoo/usms/pkg/proto/permission/v1"
	rolepb "github.com/nortoo/usms/pkg/proto/role/v1"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
)

func TestUser(t *testing.T) {
	client := pb.NewServiceClient(conn)
	roleCli := rolepb.NewServiceClient(conn)
	applicationCli := applicationpb.NewServiceClient(conn)
	menuCli := menupb.NewServiceClient(conn)
	permissionCli := permissionpb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app, err := applicationCli.Create(ctx, &applicationpb.CreateReq{
		Name:    "web",
		Comment: "",
	})
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}
	defer func() {
		_, err = applicationCli.Delete(ctx, &applicationpb.DeleteReq{Id: app.Id})
		if err != nil {
			t.Fatalf("Failed to delete application after test: %v", err)
		}
	}()

	menu, err := menuCli.Create(ctx, &menupb.CreateReq{
		ParentId: 0,
		Name:     "users",
		Path:     "/users",
		Comment:  "",
	})
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}
	defer func() {
		_, err = menuCli.Delete(ctx, &menupb.DeleteReq{Id: menu.Id})
		if err != nil {
			t.Fatalf("Failed to delete menu after test: %v", err)
		}
	}()

	permission, err := permissionCli.Create(ctx, &permissionpb.CreateReq{
		Action:   "GET",
		Resource: "/users",
		Comment:  "",
	})
	if err != nil {
		t.Fatalf("Failed to create permission: %v", err)
	}
	defer func() {
		_, err = permissionCli.Delete(ctx, &permissionpb.DeleteReq{Id: permission.Id})
		if err != nil {
			t.Fatalf("Failed to delete permission after test: %v", err)
		}
	}()

	role, err := roleCli.Create(ctx, &rolepb.CreateReq{
		Name:          "user",
		Comment:       "this is a normal user role.",
		ApplicationId: app.GetId(),
		Menus:         []uint64{menu.GetId()},
		Permissions:   []uint64{permission.GetId()},
	})
	if err != nil {
		t.Fatalf("Failed to create role: %v", err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: role.Id})
		if err != nil {
			t.Fatalf("Failed to delete role after test: %v", err)
		}
	}()

	user, err := client.Create(ctx, &pb.CreateReq{
		Username: "admin",
		Password: "hello123",
		Email:    "hello@example.com",
		Mobile:   "",
		Roles:    []int64{int64(role.Id)},
		State:    1,
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: user.Id})
		if err != nil {
			t.Fatalf("Failed to delete user after test: %v", err)
		}
	}()

	user, err = client.Update(ctx, &pb.UpdateReq{
		Id:       user.GetId(),
		Password: "abc9527",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	loginResp, err := client.Login(ctx, &pb.LoginReq{
		Identifier: "hello@example.com",
		Password:   "abc9527",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	authPass, err := client.Auth(ctx, &pb.AuthReq{
		Token:    loginResp.Token,
		Tenant:   app.Name,
		Resource: permission.GetResource(),
		Action:   permission.GetAction(),
	})
	if err != nil {
		t.Fatalf("Failed to authorize user: %v", err)
	}
	if !authPass.GetAuthorized() {
		t.Fatalf("Failed to authorize user, it should be authorized")
	}

	authDenied, err := client.Auth(ctx, &pb.AuthReq{
		Token:    loginResp.GetToken(),
		Tenant:   app.GetName(),
		Resource: "/unknown",
		Action:   "GET",
	})
	if err != nil {
		t.Fatalf("Failed to authorize user: %v", err)
	}
	if authDenied.GetAuthorized() {
		t.Fatalf("Failed to authorize user, it should not be authorized")
	}

	_, err = client.Signup(ctx, &pb.SignupReq{
		Username: "example",
		Password: "password7788",
		//Email:    "",
		//Mobile:   "",
	})
	if err != nil {
		t.Fatalf("Failed to signup: %v", err)
	}

	user, err = client.Get(ctx, &pb.GetReq{Id: user.GetId()})
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	users, err := client.List(ctx, &pb.ListReq{
		Pagination: &pbtypes.Pagination{
			Page:     1,
			PageSize: 10,
		},
	})
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users.GetItems()) != 2 {
		t.Fatalf("Failed to list users, expected at least 2 users, got %d", len(users.GetItems()))
	}
}
