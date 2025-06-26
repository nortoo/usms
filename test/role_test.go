package test

import (
	"context"
	"testing"
	"time"

	applicationpb "github.com/nortoo/usms/pkg/proto/application/v1"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	menupb "github.com/nortoo/usms/pkg/proto/menu/v1"
	permissionpb "github.com/nortoo/usms/pkg/proto/permission/v1"
	pb "github.com/nortoo/usms/pkg/proto/role/v1"
)

func TestRole(t *testing.T) {
	client := pb.NewServiceClient(conn)
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

	r, err := client.Create(ctx, &pb.CreateReq{
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
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: r.Id})
		if err != nil {
			t.Fatalf("Failed to delete role after test: %v", err)
		}
	}()

	r, err = client.Update(ctx, &pb.UpdateReq{
		Id:      r.GetId(),
		Comment: "there is a new comment.",
		//Menus:       nil,
		//Permissions: nil,
	})
	if err != nil {
		t.Fatalf("Failed to update role: %v", err)
	}
	if r.GetComment() != "there is a new comment." {
		t.Fatalf("Failed to update role, expected comment 'there is a new comment.', got '%s'", r.GetComment())
	}

	r, err = client.Get(ctx, &pb.GetReq{Id: r.Id})
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	roles, err := client.List(ctx, &pb.ListReq{Pagination: &pbtypes.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatalf("Failed to list roles: %v", err)
	}
	if len(roles.Items) == 0 {
		t.Fatalf("Failed to list roles, expected at least one role, got %d", len(roles.Items))
	}
}
