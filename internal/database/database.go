package database

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

var (
	db *gorm.DB
)

type Config struct {
	File     string `yaml:"file"`
	LogLevel string `yaml:"log_level"`
}

func Setup(cfg Config) (err error) {
	dbFile := cfg.File
	logLevel := cfg.LogLevel

	if !strings.HasSuffix(dbFile, ".db") {
		return errors.New("無效的 db_file 名稱")
	}

	db, err = gorm.Open(
		sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", dbFile)),
		&gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "[GORM]", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				LogLevel: func() logger.LogLevel {
					switch strings.ToLower(logLevel) {
					case "info":
						return logger.Info
					case "warn":
						return logger.Warn
					case "error":
						return logger.Error
					case "silent":
						return logger.Silent
					default:
						return logger.Warn
					}
				}(),
			}),
		})

	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	if db == nil {
		return nil
	}

	_db, err := db.DB()
	if err != nil {
		return err
	}

	return _db.Close()
}
