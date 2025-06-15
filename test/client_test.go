package test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

var conn *grpc.ClientConn

func TestMain(m *testing.M) {
	cert, err := tls.LoadX509KeyPair("certs/client/client.crt", "certs/client/client.key")
	if err != nil {
		fmt.Println("tls.LoadX509KeyPair err:", err)
		os.Exit(1)
	}

	certPool := x509.NewCertPool()
	rootBuf, err := os.ReadFile("certs/ca/ca.pem")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		os.Exit(1)
	}
	if !certPool.AppendCertsFromPEM(rootBuf) {
		fmt.Println("AppendCertsFromPEM err")
		os.Exit(1)
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   "USMS-SERVER",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	conn, err = grpc.NewClient("198.19.249.3:8080",
		grpc.WithTransportCredentials(creds),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   10 * time.Second,
			},
			MinConnectTimeout: 10 * time.Second,
		}))
	if err != nil {
		fmt.Println("grpc.NewClient err:", err)
		os.Exit(1)
	}

	m.Run()
}
