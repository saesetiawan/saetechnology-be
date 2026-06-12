package database

import (
	"fmt"
	"time"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewPostgresql(cfg *config.Config) *gorm.DB {
	sslMode := "disable"

	if cfg.DatabaseConfig.Secure {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Username,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Name,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: getGormLogger(cfg.AppEnv),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	// Pool configuration
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	return db
}

func getGormLogger(appEnv string) gormLogger.Interface {
	if appEnv == "production" {
		return gormLogger.Default.LogMode(gormLogger.Error)
	}

	return gormLogger.Default.LogMode(gormLogger.Info)
}
