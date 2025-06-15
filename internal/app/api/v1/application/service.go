package application

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	pb "github.com/nortoo/usms/pkg/proto/application/v1"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	"github.com/nortoo/utils-go/char"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.Application, error) {
	app := &model.Application{
		ID:        uint(snowflake.GetSnowWorker().NextId()),
		Name:      req.GetName(),
		APPID:     string(char.RandomBytes(32)),
		SecretKey: string(char.RandomBytes(64)),
		Comment:   req.GetComment(),
		State:     0,
	}
	err := _usm.Client().CreateApplication(app)
	if err != nil {
		return nil, err
	}

	return &pb.Application{
		Id:        uint64(app.ID),
		Name:      app.Name,
		Appid:     app.APPID,
		SecretKey: app.SecretKey,
		Comment:   app.Comment,
		State:     int32(app.State),
		Time: &pbtypes.TimeModel{
			CreatedAt: app.CreatedAt.Unix(),
			UpdatedAt: app.UpdatedAt.Unix(),
		},
	}, nil
}

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	return _usm.Client().DeleteApplication(&model.Application{ID: uint(req.GetId())})
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.Application, error) {
	app := &model.Application{ID: uint(req.GetId())}

	var cols []string
	if req.GetComment() != "" {
		app.Comment = req.GetComment()
		cols = append(cols, "Comment")
	}
	if req.GetState() != 0 {
		app.State = int8(req.GetState())
		cols = append(cols, "State")
	}
	if len(cols) == 0 {
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := _usm.Client().UpdateApplication(app, cols...)
	if err != nil {
		return nil, err
	}

	return &pb.Application{
		Id:        uint64(app.ID),
		Name:      app.Name,
		Appid:     app.APPID,
		SecretKey: app.SecretKey,
		Comment:   app.Comment,
		State:     int32(app.State),
		Time: &pbtypes.TimeModel{
			CreatedAt: app.CreatedAt.Unix(),
			UpdatedAt: app.UpdatedAt.Unix(),
		},
	}, nil
}

func Get(ctx context.Context, req *pb.GetReq) (*pb.Application, error) {
	app := &model.Application{
		ID:    uint(req.GetId()),
		Name:  req.GetName(),
		APPID: req.GetAppid(),
	}
	app, err := _usm.Client().GetApplication(app)
	if err != nil {
		return nil, err
	}

	return &pb.Application{
		Id:        uint64(app.ID),
		Name:      app.Name,
		Appid:     app.APPID,
		SecretKey: app.SecretKey,
		Comment:   app.Comment,
		State:     int32(app.State),
		Time: &pbtypes.TimeModel{
			CreatedAt: app.CreatedAt.Unix(),
			UpdatedAt: app.UpdatedAt.Unix(),
		},
	}, nil
}

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := _usm.Client().ListApplications(&types.QueryApplicationOptions{Pagination: &types.Pagination{
		Page:     int(req.GetPagination().Page),
		PageSize: int(req.GetPagination().GetPageSize()),
	}})
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Application, len(ret))
	for i, app := range ret {
		items[i] = &pb.Application{
			Id:        uint64(app.ID),
			Name:      app.Name,
			Appid:     app.APPID,
			SecretKey: app.SecretKey,
			Comment:   app.Comment,
			State:     int32(app.State),
			Time: &pbtypes.TimeModel{
				CreatedAt: app.CreatedAt.Unix(),
				UpdatedAt: app.UpdatedAt.Unix(),
			},
		}
	}

	return &pb.ListResp{
		Pagination: &pbtypes.PaginationResp{
			Page:  req.GetPagination().GetPage(),
			Total: total,
		},
		Items: items,
	}, nil
}
