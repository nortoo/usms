package user

import (
	"context"

	"time"

	"github.com/nortoo/usm"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/jwt"
	"github.com/nortoo/usms/internal/pkg/session"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/nortoo/usms/internal/pkg/types/user"
	"github.com/nortoo/usms/internal/pkg/utils/encryption"
	"github.com/nortoo/usms/internal/pkg/utils/identification"
	_validation "github.com/nortoo/usms/internal/pkg/validation"
	"github.com/nortoo/usms/pkg/errors"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"github.com/nortoo/utils-go/validation"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error)
	Delete(ctx context.Context, req *pb.DeleteReq) error
	Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error)
	Get(ctx context.Context, req *pb.GetReq) (*pb.User, error)
	List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error)

	Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error)
	Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error)
	Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*emptypb.Empty, error)
	ResetPassword(ctx context.Context, req *pb.ResetPasswordReq) (*emptypb.Empty, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenResp, error)

	DoesIdentifierExist(ctx context.Context, req *pb.DoesIdentifierExistReq) (*pb.DoesIdentifierExistResp, error)

	getUserFromToken(ctx context.Context, token string) (*model.User, error)
}

type service struct {
	config            *etc.Config
	env               *etc.Env
	jwt               jwt.Service
	session           session.Service
	usmCli            *usm.Client
	redisCli          *store.RedisCli
	validator         _validation.Service
	identificationSvc identification.Service
	logger            *zap.Logger
}

func NewService(
	config *etc.Config, env *etc.Env, usmCli *usm.Client, jwt jwt.Service,
	session session.Service, redisCli *store.RedisCli, validator _validation.Service,
	logger *zap.Logger) Service {
	return &service{
		config:            config,
		env:               env,
		session:           session,
		jwt:               jwt,
		usmCli:            usmCli,
		redisCli:          redisCli,
		validator:         validator,
		identificationSvc: identification.New(validator),
		logger:            logger,
	}
}

func (s *service) Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error) {
	if _, err := s.validator.IsValidUsername(req.GetUsername()); err != nil {
		return nil, errors.ErrInvalidParams.WithDetail(err.Error())
	}
	if _, err := s.validator.IsValidPassword(req.GetPassword()); err != nil {
		return nil, errors.ErrInvalidParams.WithDetail(err.Error())
	}
	if req.GetEmail() != "" {
		if !validation.IsValidEmail(req.GetEmail()) {
			return nil, errors.ErrInvalidParams.WithDetail("invalid email")
		}
	}
	if req.GetMobile() != "" {
		isValidMobile, err := validation.IsValidMobileNumber(req.GetMobile(), "US")
		if !isValidMobile || err != nil {
			return nil, errors.ErrInvalidParams.WithDetail("invalid mobile")
		}
	}

	usernameExists, err := s.usmCli.DoesUsernameExist(req.GetUsername())
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.ErrUserExists.WithDetail("username already exists")
	}
	if req.GetEmail() != "" {
		emailExists, err := s.usmCli.DoesEmailExist(req.GetEmail())
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.ErrUserExists.WithDetail("email already exists")
		}
	}
	if req.GetMobile() != "" {
		mobileExists, err := s.usmCli.DoesMobileExist(req.GetMobile())
		if err != nil {
			return nil, err
		}
		if mobileExists {
			return nil, errors.ErrUserExists.WithDetail("mobile already exists")
		}
	}

	password, err := encryption.EncryptPassword(req.GetPassword())
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	u := &model.User{
		ID:       uint(snowflake.GetSnowWorker().NextId()),
		Username: req.GetUsername(),
		Password: password,
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		State:    int8(req.GetState()),
	}
	if u.State == 0 {
		u.State = s.config.App.Settings.DefaultValue.UserState
	}

	var roles []*model.Role
	if len(req.GetRoles()) != 0 {
		for _, rid := range req.GetRoles() {
			role, err := s.usmCli.GetRole(&model.Role{ID: uint(rid)})
			if err != nil {
				s.logger.Warn("tole doesn't exist", zap.Int64("id", rid))
				continue
			}
			roles = append(roles, role)
		}
	} else {
		// assign default roles if no roles are provided.
		roles, _, err = s.usmCli.ListRoles(&types.QueryRoleOptions{
			IsDefault: []bool{true},
			WithTotal: false,
		})
		if err != nil {
			s.logger.Warn("default role doesn't exist", zap.Error(err))
		}
	}
	u.Roles = roles

	var groups []*model.Group
	if len(req.GetGroups()) != 0 {
		for _, gid := range req.GetGroups() {
			group, err := s.usmCli.GetGroup(&model.Group{ID: uint(gid)})
			if err != nil {
				s.logger.Warn("group doesn't exist", zap.Int64("id", gid))
				continue
			}
			groups = append(groups, group)
		}
	} else {
		// assign default groups if no groups are provided.
		groups, _, err = s.usmCli.ListGroups(&types.QueryGroupOptions{
			IsDefault: []bool{true},
			WithTotal: false,
		})
		if err != nil {
			s.logger.Warn("default group doesn't exist", zap.Error(err))
		}
	}
	u.Groups = groups

	err = s.usmCli.CreateUser(u)
	if err != nil {
		return nil, err
	}

	return user.ModelToPb(u), nil
}

func (s *service) Delete(ctx context.Context, req *pb.DeleteReq) error {
	return s.usmCli.DeleteUser(&model.User{ID: uint(req.GetId())})
}

func (s *service) Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	u := &model.User{ID: uint(req.GetId())}

	var cols []string
	if req.GetMobile() != "" {
		u.Mobile = req.GetMobile()
		cols = append(cols, "Mobile")
	}
	if req.GetEmail() != "" {
		u.Email = req.GetEmail()
		cols = append(cols, "Email")
	}
	if req.GetPassword() != "" {
		if _, err := s.validator.IsValidPassword(req.GetPassword()); err != nil {
			return nil, errors.ErrInvalidParams.WithDetail(err.Error())
		}

		password, err := encryption.EncryptPassword(req.GetPassword())
		if err != nil {
			return nil, errors.ErrInternalError.WithDetail(err.Error())
		}
		u.Password = password
		cols = append(cols, "Password")
	}
	if req.GetState() != 0 {
		u.State = int8(req.GetState())
		cols = append(cols, "State")
	}

	var roles []*model.Role
	for _, rid := range req.GetRoles() {
		role, err := s.usmCli.GetRole(&model.Role{ID: uint(rid)})
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}
	if len(roles) != 0 {
		u.Roles = roles
		cols = append(cols, "Roles")
	}

	var groups []*model.Group
	for _, gid := range req.GetGroups() {
		group, err := s.usmCli.GetGroup(&model.Group{ID: uint(gid)})
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}
	if len(groups) != 0 {
		u.Groups = groups
		cols = append(cols, "Groups")
	}

	if len(cols) == 0 {
		return s.Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	if u.Email != "" {
		emailExists, err := s.usmCli.DoesEmailExist(u.Email)
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.ErrUserExists.WithDetail("email already exists")
		}
	}
	if u.Mobile != "" {
		mobileExists, err := s.usmCli.DoesMobileExist(u.Mobile)
		if err != nil {
			return nil, err
		}
		if mobileExists {
			return nil, errors.ErrUserExists.WithDetail("mobile already exists")
		}
	}
	err := s.usmCli.UpdateUser(u, cols...)
	if err != nil {
		return nil, err
	}

	return user.ModelToPb(u), nil
}

func (s *service) Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	var u *model.User
	var err error

	if req.GetId() != 0 {
		u, err = s.usmCli.GetUser(&model.User{ID: uint(req.GetId())})
	} else {
		identifier := req.GetIdentifier()
		switch s.identificationSvc.Recognize(identifier) {
		case identification.Email:
			u, err = s.usmCli.GetUser(&model.User{Email: identifier})
		case identification.Mobile:
			u, err = s.usmCli.GetUser(&model.User{Mobile: identifier})
		case identification.Username:
			u, err = s.usmCli.GetUser(&model.User{Username: identifier})
		default:
			return nil, errors.ErrInvalidParams.WithDetail("invalid identifier")
		}
	}

	if err != nil {
		return nil, err
	}
	return user.ModelToPb(u), nil
}

func (s *service) List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
	options := &types.QueryUserOptions{
		Pagination: &types.Pagination{
			Page:     int(req.GetPagination().GetPage()),
			PageSize: int(req.GetPagination().GetPageSize()),
		},
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		GroupID:  uint(req.GetGroupId()),
		RoleID:   uint(req.GetRoleId()),
	}
	for _, state := range req.GetState() {
		options.States = append(options.States, int8(state))
	}
	ret, total, err := s.usmCli.ListUsers(options)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.User, len(ret))
	for i, r := range ret {
		items[i] = user.ModelToPb(r)
	}

	return &pb.ListResp{
		Pagination: &pbtypes.PaginationResp{
			Page:  req.GetPagination().Page,
			Total: total,
		},
		Items: items,
	}, nil
}

func (s *service) Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error) {
	_, err := s.Create(ctx, &pb.CreateReq{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		State:    int32(s.config.App.Settings.DefaultValue.UserState),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *service) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	var u *model.User
	var err error

	identifier := req.GetIdentifier()
	switch s.identificationSvc.Recognize(identifier) {
	case identification.Email:
		u, err = s.usmCli.GetUser(&model.User{Email: identifier})
	case identification.Mobile:
		u, err = s.usmCli.GetUser(&model.User{Mobile: identifier})
	case identification.Username:
		u, err = s.usmCli.GetUser(&model.User{Username: identifier})
	default:
		return nil, errors.ErrInvalidParams.WithDetail("invalid identifier")
	}
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}
	if !encryption.ComparePassword(u.Password, req.GetPassword()) {
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}

	accessToken, refreshToken, err := s.jwt.IssueAccessTokenAndRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResp{
		User:         user.ModelToPb(u),
		Token:        accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *service) getUserFromToken(ctx context.Context, token string) (*model.User, error) {
	claims, err := s.jwt.ParseToken(token, s.env.JWTSecretKey)
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	tokenId := claims.ID
	if tokenId == "" {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	// check if the token is in the session blocklist.
	sessionBlocklistKey := s.session.GenerateSessionBlocklistKey(tokenId)
	exist, err := s.redisCli.GetRDB().Exists(ctx, sessionBlocklistKey).Result()
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}
	if exist > 0 {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	userID := claims.UID
	if userID <= 0 {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid user ID in token")
	}

	u, err := s.usmCli.GetUser(&model.User{
		ID:    uint(userID),
		Roles: []*model.Role{},
	})
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("user not found")
	}

	return u, nil
}

func (s *service) Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
	u, err := s.getUserFromToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}
	if u.Roles == nil {
		return &pb.AuthResp{Authorized: false}, nil
	}

	for _, role := range u.Roles {
		authorized, err := s.usmCli.Authorize(role.Name, req.GetTenant(), req.GetResource(), req.GetAction())
		if err != nil {
			s.logger.Warn("authorize fail", zap.Error(err))
			continue
		}
		if authorized {
			return &pb.AuthResp{
				Uid:        uint64(u.ID),
				Authorized: true,
			}, nil
		}
	}
	return &pb.AuthResp{Authorized: false}, nil
}

func (s *service) DoesIdentifierExist(ctx context.Context, req *pb.DoesIdentifierExistReq) (*pb.DoesIdentifierExistResp, error) {
	var resp pb.DoesIdentifierExistResp

	if req.GetEmail() != "" {
		emailExists, err := s.usmCli.DoesEmailExist(req.GetEmail())
		if err != nil {
			return nil, err
		}
		resp.EmailExist = emailExists
	}
	if req.GetMobile() != "" {
		mobileExists, err := s.usmCli.DoesMobileExist(req.GetMobile())
		if err != nil {
			return nil, err
		}
		resp.MobileExist = mobileExists
	}
	if req.GetUsername() != "" {
		usernameExists, err := s.usmCli.DoesUsernameExist(req.GetUsername())
		if err != nil {
			return nil, err
		}
		resp.UsernameExist = usernameExists
	}

	return &resp, nil
}

func (s *service) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*emptypb.Empty, error) {
	u, err := s.usmCli.GetUser(&model.User{ID: uint(req.GetUid())})
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.GetOldPassword())); err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}

	if _, err := s.validator.IsValidPassword(req.GetNewPassword()); err != nil {
		return nil, errors.ErrInvalidParams.WithDetail(err.Error())
	}

	password, err := encryption.EncryptPassword(req.GetNewPassword())
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}
	u.Password = password

	err = s.usmCli.UpdateUser(u, "Password")
	if err != nil {
		return nil, err
	}
	s.session.RevokeUserTokens(ctx, u.ID)

	return &emptypb.Empty{}, nil
}

func (s *service) ResetPassword(ctx context.Context, req *pb.ResetPasswordReq) (*emptypb.Empty, error) {
	u, err := s.usmCli.GetUser(&model.User{ID: uint(req.GetUid())})
	if err != nil {
		return nil, err
	}

	if _, err := s.validator.IsValidPassword(req.GetNewPassword()); err != nil {
		return nil, errors.ErrInvalidParams.WithDetail(err.Error())
	}

	password, err := encryption.EncryptPassword(req.GetNewPassword())
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}
	u.Password = password

	err = s.usmCli.UpdateUser(u, "Password")
	if err != nil {
		return nil, err
	}
	s.session.RevokeUserTokens(ctx, u.ID)

	return &emptypb.Empty{}, nil
}

func (s *service) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenResp, error) {
	claims, err := s.jwt.ParseToken(req.GetRefreshToken(), s.env.JWTRefreshSecretKey)
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token.")
	}

	tokenId := claims.ID
	if tokenId == "" {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	refreshTokenBlocklistKey := s.session.GenerateSessionRefreshTokenBlocklistStoreKey(tokenId)
	exist, err := s.redisCli.GetRDB().Exists(ctx, refreshTokenBlocklistKey).Result()
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}
	if exist > 0 {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	userID := claims.UID
	if userID <= 0 {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid user ID in token")
	}

	accessToken, refreshToken, err := s.jwt.IssueAccessTokenAndRefreshToken(uint(userID))
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token.")
	}

	// revoke the refresh token once it is used.
	// when get the expiresAt failed, use the default expire of the token.
	refreshTokenStoreKey := s.session.GenerateSessionRefreshTokenStoreKey(uint(userID), tokenId)
	expire := time.Duration(s.config.App.Settings.JWT.TokenRefreshTime) * time.Second
	ttl, err := s.redisCli.GetRDB().TTL(ctx, refreshTokenStoreKey).Result()
	if err == nil && ttl > 5 {
		// when ttl is less than or equal 5, we consider that this token is about to expire.
		expire = ttl
	}

	err = s.redisCli.GetRDB().Set(ctx, refreshTokenBlocklistKey, "", expire).Err()
	if err != nil {
		s.logger.Warn("Failed to add tokenId to blocklist", zap.Error(err))
	}

	return &pb.RefreshTokenResp{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
