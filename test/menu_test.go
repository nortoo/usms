package test

import (
	"context"
	"testing"
	"time"

	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/menu/v1"
)

func TestMenu(t *testing.T) {
	client := pb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := client.Create(ctx, &pb.CreateReq{
		ParentId: 0,
		Name:     "resource",
		Path:     "/resource",
		Comment:  "the root directory of resources",
	})
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: m.Id})
		if err != nil {
			t.Fatalf("Failed to delete menu: %v", err)
		}
	}()

	m, err = client.Update(ctx, &pb.UpdateReq{
		Id:      m.Id,
		Name:    "resource-updated",
		Path:    "/resource-updated",
		Comment: "",
	})
	if err != nil {
		t.Fatalf("Failed to update menu: %v", err)
	}
	if m.Name != "resource-updated" {
		t.Fatalf("Failed to update menu, expected name 'resource-updated', got '%s'", m.Name)
	}

	m, err = client.Get(ctx, &pb.GetReq{Id: m.Id})
	if err != nil {
		t.Fatalf("Failed to get menu: %v", err)
	}

	menus, err := client.List(ctx, &pb.ListReq{Pagination: &pbtypes.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatalf("Failed to list menus: %v", err)
	}
	if len(menus.Items) == 0 {
		t.Fatalf("Expected to find at least one menu, got %d", len(menus.Items))
	}
}
