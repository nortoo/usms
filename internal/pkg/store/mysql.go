package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nortoo/usms/internal/pkg/etc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetStore(name DBName) *gorm.DB {
	return dbs[name]
}

func InitMysql(config *etc.Store) error {
	for name, cfg := range config.Mysql {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DB)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			SkipDefaultTransaction:                   true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Warn, // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      true,        // Don't include params in the SQL log
					Colorful:                  false,       // Disable color
				},
			),
		})
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(10 * time.Minute)

		dbs[DBName(name)] = db
	}

	return nil
}
