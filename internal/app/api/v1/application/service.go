package application

import (
	"context"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	pb "github.com/nortoo/usms/pkg/proto/application/v1"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	"github.com/nortoo/utils-go/char"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.Application, error)
	Delete(ctx context.Context, req *pb.DeleteReq) error
	Update(ctx context.Context, req *pb.UpdateReq) (*pb.Application, error)
	Get(ctx context.Context, req *pb.GetReq) (*pb.Application, error)
	List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error)
}

type service struct {
	usmCli *usm.Client
}

func NewService(usmCli *usm.Client) Service {
	return &service{usmCli: usmCli}
}

func (s *service) Create(ctx context.Context, req *pb.CreateReq) (*pb.Application, error) {
	// Todo: verify if the application name already exists
	app := &model.Application{
		ID:        uint(snowflake.GetSnowWorker().NextId()),
		Name:      req.GetName(),
		APPID:     string(char.RandomBytes(32)),
		SecretKey: string(char.RandomBytes(64)),
		Comment:   req.GetComment(),
		State:     0,
	}
	err := s.usmCli.CreateApplication(app)
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

func (s *service) Delete(ctx context.Context, req *pb.DeleteReq) error {
	return s.usmCli.DeleteApplication(&model.Application{ID: uint(req.GetId())})
}

func (s *service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.Application, error) {
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
		return s.Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	err := s.usmCli.UpdateApplication(app, cols...)
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

func (s *service) Get(ctx context.Context, req *pb.GetReq) (*pb.Application, error) {
	app := &model.Application{
		ID:    uint(req.GetId()),
		Name:  req.GetName(),
		APPID: req.GetAppid(),
	}
	app, err := s.usmCli.GetApplication(app)
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

func (s *service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	ret, total, err := s.usmCli.ListApplications(&types.QueryApplicationOptions{Pagination: &types.Pagination{
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
