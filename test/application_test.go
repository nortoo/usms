package test

import (
	"context"
	"testing"
	"time"

	pb "github.com/nortoo/usms/pkg/proto/application/v1"
)

func TestApplication(t *testing.T) {
	client := pb.NewServiceClient(conn)

	// Create a context with timeout for this specific RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Create(ctx, &pb.CreateReq{
		Name:    "Hellp-2",
		Comment: "This is a test application",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp:%v", resp)
}
