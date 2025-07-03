package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/nortoo/usms/internal/app/api"
	"github.com/nortoo/usms/internal/app/api/v1/application"
	"github.com/nortoo/usms/internal/app/api/v1/group"
	"github.com/nortoo/usms/internal/app/api/v1/menu"
	"github.com/nortoo/usms/internal/app/api/v1/permission"
	"github.com/nortoo/usms/internal/app/api/v1/role"
	"github.com/nortoo/usms/internal/app/api/v1/user"
	"github.com/nortoo/usms/internal/app/api/v1/verification"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	grpcServer *grpc.Server
	container  *api.Container
}

func NewServer(container *api.Container) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(container.Config.App.Certs.CertFile, container.Config.App.Certs.KeyFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	rootBuf, err := os.ReadFile(container.Config.App.Certs.CAFile)
	if err != nil {
		return nil, err
	}
	if !certPool.AppendCertsFromPEM(rootBuf) {
		return nil, errors.New("failed to set ca cert")
	}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	})))

	application.Register(grpcServer, container.ApplicationHandler)
	group.Register(grpcServer, container.GroupHandler)
	menu.Register(grpcServer, container.MenuHandler)
	permission.Register(grpcServer, container.PermissionHandler)
	role.Register(grpcServer, container.RoleHandler)
	user.Register(grpcServer, container.UserHandler)
	verification.Register(grpcServer, container.VerificationHandler)

	return &Server{
		grpcServer: grpcServer,
		container:  container,
	}, nil
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.container.Logger.Info(fmt.Sprintf("Starting gRPC server on :%d\n", port))
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
