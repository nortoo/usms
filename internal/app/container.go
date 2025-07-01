package app

import (
	"github.com/nortoo/usm"
	"github.com/nortoo/usms/internal/app/api/v1/application"
	"github.com/nortoo/usms/internal/app/api/v1/group"
	"github.com/nortoo/usms/internal/app/api/v1/menu"
	"github.com/nortoo/usms/internal/app/api/v1/permission"
	"github.com/nortoo/usms/internal/app/api/v1/role"
	"github.com/nortoo/usms/internal/app/api/v1/user"
	"github.com/nortoo/usms/internal/app/api/v1/verification"
	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/jwt"
	"github.com/nortoo/usms/internal/pkg/session"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/nortoo/usms/internal/pkg/utils/identification"
	_validation "github.com/nortoo/usms/internal/pkg/validation"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Container struct {
	Config *etc.Config
	Logger *zap.Logger

	ApplicationHandler  *application.Handler
	GroupHandler        *group.Handler
	MenuHandler         *menu.Handler
	PermissionHandler   *permission.Handler
	RoleHandler         *role.Handler
	UserHandler         *user.Handler
	VerificationHandler *verification.Handler
}

func NewContainer(
	config *etc.Config,
	env *etc.Env,
	logger *zap.Logger,
	casbinPolicyPath string,
) (*Container, error) {
	if err := store.InitMysql(config.Store); err != nil {
		return nil, errors.Errorf("failed to init mysql: %v\n", err)
	}
	redisCli, err := store.NewRedisCli(config.Store)
	if err != nil {
		return nil, errors.Errorf("failed to init redis: %v\n", err)
	}

	usmCli, err := usm.New(&usm.Options{
		Store: store.GetStore(store.Default),
		CasbinOptions: &usm.CasbinOptions{
			Store:      store.GetStore(store.Default),
			PolicyPath: casbinPolicyPath,
		},
	})
	if err != nil {
		logger.Fatal("Failed to init usm", zap.Error(err))
		return nil, errors.Errorf("failed to init usm: %v\n", err)
	}

	applicationHandler := application.NewHandler(application.NewService(usmCli))
	groupHandler := group.NewHandler(group.NewService(usmCli))
	menuHandler := menu.NewHandler(menu.NewService(usmCli))
	permissionHandler := permission.NewHandler(permission.NewService(usmCli))
	roleHandler := role.NewHandler(role.NewService(usmCli))

	sessionService := session.NewService(config, redisCli, logger)
	jwtService := jwt.NewService(config, env, sessionService, redisCli, logger)
	validator := _validation.New(config)
	userService := user.NewService(
		config,
		env,
		usmCli,
		jwtService,
		sessionService,
		redisCli,
		validator,
		logger,
	)
	userHandler := user.NewHandler(userService)
	verificationHandler := verification.NewHandler(verification.NewService(usmCli, identification.New(validator)))

	return &Container{
		Config: config,
		Logger: logger,

		ApplicationHandler:  applicationHandler,
		GroupHandler:        groupHandler,
		MenuHandler:         menuHandler,
		PermissionHandler:   permissionHandler,
		RoleHandler:         roleHandler,
		UserHandler:         userHandler,
		VerificationHandler: verificationHandler,
	}, nil
}
