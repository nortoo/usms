package main

import (
	"fmt"
	"os"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/nortoo/usms/internal/pkg/log"
	"github.com/nortoo/usms/internal/pkg/snowflake"
	"github.com/nortoo/usms/internal/pkg/store"
	"github.com/nortoo/usms/internal/pkg/usm"
	"go.uber.org/zap"
)

func main() {
	err := log.InitLogger("conf/log.json")
	if err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	if err = etc.Load("conf/app.yml"); err != nil {
		log.GetLogger().Fatal("Failed to load config file", zap.Error(err))
		os.Exit(1)
	}
	if err = etc.LoadEnv(); err != nil {
		log.GetLogger().Fatal("Failed to load environment variables", zap.Error(err))
		os.Exit(1)
	}
	if err = snowflake.Init(etc.GetConfig().App.SnowflakeID); err != nil {
		log.GetLogger().Fatal("Failed to init snowflake", zap.Error(err))
		os.Exit(1)
	}
	if err = store.InitMysql(etc.GetConfig().Store); err != nil {
		log.GetLogger().Fatal("Failed to init mysql", zap.Error(err))
		os.Exit(1)
	}
	if err = usm.Init(
		store.GetStore(store.Default),
		store.GetStore(store.Default),
		"conf/casbin.conf"); err != nil {
		log.GetLogger().Fatal("Failed to init usm", zap.Error(err))
		os.Exit(1)
	}
	if err = run(8080,
		etc.GetConfig().App.Certs.CertFile,
		etc.GetConfig().App.Certs.KeyFile,
		etc.GetConfig().App.Certs.CAFile,
	); err != nil {
		log.GetLogger().Fatal("Failed to run usms", zap.Error(err))
		os.Exit(1)
	}
}
