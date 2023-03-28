package commom

import (
	log2 "log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	log := logger.Warn
	if DefaultConfig.App.Mode == "debug" {
		log = logger.Info
	}
	if db, err := gorm.Open(sqlite.Open(file+"?mode=wal"), &gorm.Config{
		Logger: logger.Default.LogMode(log),
	}); err != nil {
		if err != nil {
			panic("failed to connect database")
		}
	} else {
		DB = db
		if sqlDB, err := db.DB(); err != nil {
			log2.Fatal(err)
		} else {
			sqlDB.SetMaxOpenConns(1)
		}
	}

}
