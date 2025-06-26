package test

import (
	"context"
	"testing"
	"time"

	pb "github.com/nortoo/usms/pkg/proto/application/v1"
	"github.com/nortoo/usms/pkg/proto/common/v1/types"
)

func TestApplication(t *testing.T) {
	client := pb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Create(ctx, &pb.CreateReq{
		Name:    "Hellp",
		Comment: "This is a test application",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err = client.Delete(ctx, &pb.DeleteReq{Id: resp.Id})
		if err != nil {
			t.Fatalf("Failed to delete application after test: %v", err)
		}
	}()

	app, err := client.Update(ctx, &pb.UpdateReq{
		Id:      resp.GetId(),
		Comment: "I changed the comment.",
		State:   1,
	})
	if err != nil {
		t.Fatalf("Failed to update application: %v", err)
	}
	if app.GetComment() != "I changed the comment." || app.GetState() != 1 {
		t.Fatalf("Failed to update application, expected comment 'I changed the comment.' and state 1, got comment '%s' and state %d", app.GetComment(), app.GetState())
	}

	app, err = client.Get(ctx, &pb.GetReq{Name: "Hellp"})
	if err != nil {
		t.Fatalf("Failed to get application: %v", err)
	}

	apps, err := client.List(ctx, &pb.ListReq{Pagination: &types.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatalf("Failed to list applications: %v", err)
	}
	if len(apps.Items) == 0 {
		t.Fatalf("Failed to list applications, expected at least 1 application, got %d", len(apps.Items))
	}
}
