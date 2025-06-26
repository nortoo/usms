package test

import (
	"context"
	"testing"
	"time"

	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/permission/v1"
)

func TestPermission(t *testing.T) {
	client := pb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := client.Create(ctx, &pb.CreateReq{
		Action:   "POST",
		Resource: "/api/v1/users",
		Comment:  "Create a user",
	})
	if err != nil {
		t.Fatalf("Failed to create permission: %v", err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: p.Id})
		if err != nil {
			t.Fatalf("Failed to delete permission after test: %v", err)
		}
	}()

	p, err = client.Update(ctx, &pb.UpdateReq{
		Id:      p.Id,
		Comment: "I changed the comment.",
	})
	if err != nil {
		t.Fatalf("Failed to update permission: %v", err)
	}
	if p.GetComment() != "I changed the comment." {
		t.Fatalf("Failed to update permission after update: %v", p.GetComment())
	}

	p, err = client.Get(ctx, &pb.GetReq{Id: p.Id})
	if err != nil {
		t.Fatalf("Failed to get permission: %v", err)
	}

	permissions, err := client.List(ctx, &pb.ListReq{Pagination: &pbtypes.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatalf("Failed to list permissions: %v", err)
	}
	if len(permissions.Items) == 0 {
		t.Fatalf("Failed to list permissions, expected at least one permission, got %d", len(permissions.Items))
	}
}
