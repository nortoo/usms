package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/nortoo/usms/internal/app/api/v1/application"
	"github.com/nortoo/usms/internal/app/api/v1/group"
	"github.com/nortoo/usms/internal/app/api/v1/menu"
	"github.com/nortoo/usms/internal/app/api/v1/permission"
	"github.com/nortoo/usms/internal/app/api/v1/role"
	"github.com/nortoo/usms/internal/app/api/v1/user"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func run(port int, certFile, keyFile, caFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	rootBuf, err := os.ReadFile(caFile)
	if err != nil {
		return err
	}
	if !certPool.AppendCertsFromPEM(rootBuf) {
		return errors.New("failed to set ca cert")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	})))

	application.Register(grpcServer)
	group.Register(grpcServer)
	menu.Register(grpcServer)
	permission.Register(grpcServer)
	role.Register(grpcServer)
	user.Register(grpcServer)

	fmt.Printf("gRPC server is running on :%d\n", port)
	return grpcServer.Serve(lis)
}
