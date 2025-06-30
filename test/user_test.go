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
		_, err = roleCli.Delete(ctx, &rolepb.DeleteReq{Id: role.Id})
		if err != nil {
			t.Fatalf("Failed to delete role after test: %v", err)
		}
	}()

	var userCount int
	user, err := client.Create(ctx, &pb.CreateReq{
		Username: "admin",
		Password: "Hello123@",
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
	userCount += 1

	_, err = client.Update(ctx, &pb.UpdateReq{
		Id:       user.GetId(),
		Password: "Abc95279527.",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	_, err = client.Get(ctx, &pb.GetReq{Id: user.GetId()})
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	loginResp, err := client.Login(ctx, &pb.LoginReq{
		Identifier: "hello@example.com",
		Password:   "Abc95279527.",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Passed Auth.
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

	// Denied Auth.
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

	type signupSample struct {
		Request *pb.SignupReq
		Allowed bool
	}
	signupSamples := []*signupSample{
		{
			Request: &pb.SignupReq{
				Username: "example",
				Password: "Password7788.",
				Email:    "random-1@example.com",
				Mobile:   "+15183098999",
			},
			Allowed: true,
		},
		{
			Request: &pb.SignupReq{
				// duplicated username
				Username: "example",
				Password: "Password7788.",
				Email:    "random-2@example.com",
			},
			Allowed: false,
		},
		{
			Request: &pb.SignupReq{
				Username: "example-2",

				// invalid password
				Password: "password7788.",
				Email:    "random-3@example.com",
			},
			Allowed: false,
		},
		{
			Request: &pb.SignupReq{
				Username: "example-3",
				Password: "Password7788.",

				// duplicated email
				Email: "random-3@example.com",
			},
			Allowed: false,
		},
	}

	for _, s := range signupSamples {
		_, err = client.Signup(ctx, s.Request)
		if s.Allowed && err != nil {
			t.Fatalf("Failed to signup: %v", err)
		}
		if !s.Allowed && err == nil {
			t.Fatalf("The duplicated email should be rejectedd, but it passed.")
		}
	}
	userCount += 1

	identifiers := []string{
		signupSamples[0].Request.Username,
		signupSamples[0].Request.Email,
		signupSamples[0].Request.Mobile,
	}
	password := signupSamples[0].Request.Password
	for _, identifier := range identifiers {
		t.Logf("Testing identifier %s", identifier)
		_, err := client.Login(ctx, &pb.LoginReq{
			Identifier: identifier,
			Password:   password,
		})
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
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

	if len(users.GetItems()) != userCount {
		t.Fatalf("Failed to list users, expected %d users, got %d", userCount, len(users.GetItems()))
	}

	newPassword := "NewPassword123.."
	for _, u := range users.GetItems() {
		_, err = client.ResetPassword(ctx, &pb.ResetPasswordReq{
			Uid:         u.GetId(),
			NewPassword: newPassword,
		})
		if err != nil {
			t.Fatalf("Failed to reset password: %v", err)
		}

		_, err = client.Login(ctx, &pb.LoginReq{
			Identifier: u.GetUsername(),
			Password:   newPassword,
		})
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
	}

}
