package test

import (
	"context"
	"testing"
	"time"

	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/usergroup/v1"
)

func TestGroup(t *testing.T) {
	client := pb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	g, err := client.Create(ctx, &pb.CreateReq{
		Name:    "admin",
		Comment: "this is a group for admin users",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: g.Id})
		if err != nil {
			t.Fatalf("Failed to delete group after test: %v", err)
		}
	}()

	g, err = client.Update(ctx, &pb.UpdateReq{
		Id:      g.Id,
		Name:    "admin-group",
		Comment: "this is a group for admin users, updated",
	})
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}
	if g.Name != "admin-group" {
		t.Fatalf("Failed to update group, expected name 'admin-group', got '%s'", g.Name)
	}

	g, err = client.Get(ctx, &pb.GetReq{Id: g.Id})
	if err != nil {
		t.Fatalf("Failed to get group: %v", err)
	}

	groups, err := client.List(ctx, &pb.ListReq{Pagination: &pbtypes.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatalf("Failed to list groups: %v", err)
	}
	if len(groups.Items) == 0 {
		t.Fatalf("Failed to list groups, expected at least one group, got %d", len(groups.Items))
	}
}
