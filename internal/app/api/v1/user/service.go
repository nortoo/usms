package user

import (
	"context"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/jwt"
	"github.com/nortoo/usms/internal/pkg/log"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/internal/pkg/types/user"
	_usm "github.com/nortoo/usms/internal/pkg/usm"
	_validation "github.com/nortoo/usms/internal/pkg/validation"
	"github.com/nortoo/usms/pkg/errors"
	pbtypes "github.com/nortoo/usms/pkg/proto/common/v1/types"
	pb "github.com/nortoo/usms/pkg/proto/user/v1"
	"github.com/nortoo/utils-go/validation"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func Create(ctx context.Context, req *pb.CreateReq) (*pb.User, error) {
	if !_validation.IsValidUsername(req.GetUsername()) {
		return nil, errors.ErrInvalidParams.WithDetail("invalid username")
	}
	if !_validation.IsValidPassword(req.GetPassword()) {
		return nil, errors.ErrInvalidParams.WithDetail("invalid password")
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

	usernameExists, err := _usm.Client().DoesUsernameExist(req.GetUsername())
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.ErrUserExists.WithDetail("username already exists")
	}
	if req.GetEmail() != "" {
		emailExists, err := _usm.Client().DoesEmailExist(req.GetEmail())
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.ErrUserExists.WithDetail("email already exists")
		}
	}
	if req.GetMobile() != "" {
		mobileExists, err := _usm.Client().DoesMobileExist(req.GetMobile())
		if err != nil {
			return nil, err
		}
		if mobileExists {
			return nil, errors.ErrUserExists.WithDetail("mobile already exists")
		}
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	u := &model.User{
		Model:    gorm.Model{ID: uint(snowflake.GetSnowWorker().NextId())},
		Username: req.GetUsername(),
		Password: string(password),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		State:    int8(req.GetState()),
	}
	if u.State == 0 {
		u.State = etc.GetConfig().App.Settings.DefaultValue.UserState
	}

	var roles []*model.Role
	if len(req.GetRoles()) != 0 {
		for _, rid := range req.GetRoles() {
			role, err := _usm.Client().GetRole(&model.Role{ID: uint(rid)})
			if err != nil {
				log.GetLogger().Warn("tole doesn't exist", zap.Int64("id", rid))
				continue
			}
			roles = append(roles, role)
		}
	} else {
		// assign default roles if no roles are provided.
		roles, _, err = _usm.Client().ListRoles(&types.QueryRoleOptions{
			IsDefault: []bool{true},
			WithTotal: false,
		})
		if err != nil {
			log.GetLogger().Warn("default role doesn't exist", zap.Error(err))
		}
	}
	u.Roles = roles

	var groups []*model.Group
	if len(req.GetGroups()) != 0 {
		for _, gid := range req.GetGroups() {
			group, err := _usm.Client().GetGroup(&model.Group{ID: uint(gid)})
			if err != nil {
				log.GetLogger().Warn("group doesn't exist", zap.Int64("id", gid))
				continue
			}
			groups = append(groups, group)
		}
	} else {
		// assign default groups if no groups are provided.
		groups, _, err = _usm.Client().ListGroups(&types.QueryGroupOptions{
			IsDefault: []bool{true},
			WithTotal: false,
		})
		if err != nil {
			log.GetLogger().Warn("default group doesn't exist", zap.Error(err))
		}
	}
	u.Groups = groups

	err = _usm.Client().CreateUser(u)
	if err != nil {
		return nil, err
	}

	return user.ModelToPb(u), nil
}

func Delete(ctx context.Context, req *pb.DeleteReq) error {
	return _usm.Client().DeleteUser(&model.User{Model: gorm.Model{ID: uint(req.GetId())}})
}

func Update(ctx context.Context, req *pb.UpdateReq) (*pb.User, error) {
	u := &model.User{Model: gorm.Model{ID: uint(req.GetId())}}

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
		if !_validation.IsValidPassword(req.GetPassword()) {
			return nil, errors.ErrInvalidParams.WithDetail("invalid password")
		}

		password, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.ErrInternalError.WithDetail(err.Error())
		}
		u.Password = string(password)
		cols = append(cols, "Password")
	}
	if req.GetState() != 0 {
		u.State = int8(req.GetState())
		cols = append(cols, "State")
	}

	var roles []*model.Role
	for _, rid := range req.GetRoles() {
		role, err := _usm.Client().GetRole(&model.Role{ID: uint(rid)})
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
		group, err := _usm.Client().GetGroup(&model.Group{ID: uint(gid)})
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
		return Get(ctx, &pb.GetReq{Id: req.GetId()})
	}

	if u.Email != "" {
		emailExists, err := _usm.Client().DoesEmailExist(u.Email)
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.ErrUserExists.WithDetail("email already exists")
		}
	}
	if u.Mobile != "" {
		mobileExists, err := _usm.Client().DoesMobileExist(u.Mobile)
		if err != nil {
			return nil, err
		}
		if mobileExists {
			return nil, errors.ErrUserExists.WithDetail("mobile already exists")
		}
	}
	err := _usm.Client().UpdateUser(u, cols...)
	if err != nil {
		return nil, err
	}

	return user.ModelToPb(u), nil
}

func Get(ctx context.Context, req *pb.GetReq) (*pb.User, error) {
	u, err := _usm.Client().GetUser(&model.User{Model: gorm.Model{ID: uint(req.GetId())}})
	if err != nil {
		return nil, err
	}
	return user.ModelToPb(u), nil
}

func List(ctx context.Context, req *pb.ListReq) (*pb.ListResp, error) {
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
	ret, total, err := _usm.Client().ListUsers(options)
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

func Signup(ctx context.Context, req *pb.SignupReq) (*emptypb.Empty, error) {
	_, err := Create(ctx, &pb.CreateReq{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Mobile:   req.GetMobile(),
		State:    int32(etc.GetConfig().App.Settings.DefaultValue.UserState),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// loginIdentifier represents the identity type for login including username, email, and mobile.
type loginIdentifier string

const (
	identifierUsername loginIdentifier = "username"
	identifierEmail    loginIdentifier = "email"
	identifierMobile   loginIdentifier = "mobile"
)

func recognizeLoginIdentity(req *pb.LoginReq) loginIdentifier {
	identifier := req.GetIdentifier()

	if validation.IsValidEmail(identifier) {
		return identifierEmail
	} else if isValid, _ := validation.IsValidMobileNumber(identifier, "US"); isValid {
		return identifierMobile
	} else if _validation.IsValidUsername(identifier) {
		return identifierUsername
	} else {
		return "unknown"
	}
}

func Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	var u *model.User
	var err error

	identifier := req.GetIdentifier()
	switch recognizeLoginIdentity(req) {
	case identifierEmail:
		u, err = _usm.Client().GetUser(&model.User{Email: identifier})
	case identifierMobile:
		u, err = _usm.Client().GetUser(&model.User{Mobile: identifier})
	case identifierUsername:
		u, err = _usm.Client().GetUser(&model.User{Username: identifier})
	default:
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.GetPassword())); err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("username or password is incorrect")
	}

	token, err := jwt.GenerateToken(u.ID, etc.GetEnv().JWTSecretKey, etc.GetConfig().App.Settings.JWT.TokenExpireTime)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	refreshToken, err := jwt.GenerateToken(
		u.ID,
		etc.GetEnv().JWTRefreshSecretKey,
		etc.GetConfig().App.Settings.JWT.TokenRefreshTime)
	if err != nil {
		return nil, errors.ErrInternalError.WithDetail(err.Error())
	}

	return &pb.LoginResp{
		User:         user.ModelToPb(u),
		Token:        token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func getUserFromToken(ctx context.Context, token string) (*model.User, error) {
	claims, err := jwt.ParseToken(token, etc.GetEnv().JWTSecretKey)
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid token")
	}

	userID := claims.UID
	if userID <= 0 {
		return nil, errors.ErrUnauthenticated.WithDetail("invalid user ID in token")
	}

	u, err := _usm.Client().GetUser(&model.User{
		Model: gorm.Model{ID: uint(userID)},
		Roles: []*model.Role{},
	})
	if err != nil {
		return nil, errors.ErrUnauthenticated.WithDetail("user not found")
	}

	return u, nil
}

func Auth(ctx context.Context, req *pb.AuthReq) (*pb.AuthResp, error) {
	u, err := getUserFromToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}
	if u.Roles == nil {
		return &pb.AuthResp{Authorized: false}, nil
	}

	for _, role := range u.Roles {
		authorized, err := _usm.Client().Authorize(role.Name, req.GetTenant(), req.GetResource(), req.GetAction())
		if err != nil {
			log.GetLogger().Warn("authorize fail", zap.Error(err))
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

func DoesIdentifierExist(ctx context.Context, req *pb.DoesIdentifierExistReq) (*pb.DoesIdentifierExistResp, error) {
	var resp pb.DoesIdentifierExistResp

	if req.GetEmail() != "" {
		emailExists, err := _usm.Client().DoesEmailExist(req.GetEmail())
		if err != nil {
			return nil, err
		}
		resp.EmailExist = emailExists
	}
	if req.GetMobile() != "" {
		mobileExists, err := _usm.Client().DoesMobileExist(req.GetMobile())
		if err != nil {
			return nil, err
		}
		resp.MobileExist = mobileExists
	}
	if req.GetUsername() != "" {
		usernameExists, err := _usm.Client().DoesUsernameExist(req.GetUsername())
		if err != nil {
			return nil, err
		}
		resp.UsernameExist = usernameExists
	}

	return &resp, nil
}
