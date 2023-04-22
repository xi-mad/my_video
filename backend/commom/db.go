package commom

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var DB *gorm.DB

func Load(file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if absPath, err := filepath.Abs(file); err != nil {
			panic(err)
		} else {
			if f, err := os.Create(absPath); err != nil {
				panic(err)
			} else {
				_ = f.Close()
			}
		}
	}

	if db, err := gorm.Open(sqlite.Open(file+"?mode=wal"), &gorm.Config{
		Logger: logger.New(log.New(io.MultiWriter(Logfile, os.Stdout), "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel(),
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		}),
	}); err != nil {
		if err != nil {
			panic("failed to connect database")
		}
	} else {
		DB = db
		if sqlDB, err := db.DB(); err != nil {
			log.Fatal(err)
		} else {
			sqlDB.SetMaxOpenConns(1)
		}
	}

}

func logLevel() logger.LogLevel {
	logLevel := logger.Warn
	if DefaultConfig.App.Mode == "debug" {
		logLevel = logger.Info
	}
	return logLevel
}
