package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nortoo/logger"
	"github.com/nortoo/usms/internal/app/api"
	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/grpc"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"go.uber.org/zap"
)

func main() {
	logCfgPath := flag.String("logcfg", "conf/log.json", "logger config file path")
	configPath := flag.String("c", "conf/app.yml", "config file path")

	flag.Parse()

	if logCfgPath == nil || *logCfgPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	log, err := logger.New(*logCfgPath)
	if err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	config, err := etc.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config file", zap.Error(err))
		os.Exit(1)
	}
	envCfg, err := etc.LoadEnv()
	if err != nil {
		log.Fatal("Failed to load environment variables", zap.Error(err))
		os.Exit(1)
	}

	if err = snowflake.Init(config.App.SnowflakeID); err != nil {
		log.Fatal("Failed to init snowflake", zap.Error(err))
		os.Exit(1)
	}

	container, err := api.NewContainer(config, envCfg, log, "conf/casbin.conf")
	if err != nil {
		log.Fatal("Failed to init container", zap.Error(err))
		os.Exit(1)
	}

	s, err := grpc.NewServer(container)
	if err != nil {
		log.Fatal("Failed to init server", zap.Error(err))
		os.Exit(1)
	}

	go func() {
		if err := s.Start(8080); err != nil {
			log.Fatal("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down server...")

	// Graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Stop()

	container.Logger.Info("Server stopped")
}
