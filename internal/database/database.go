package database

import (
	"errors"
	"fmt"
	"github.com/Wuchieh/candy-house-bot/internal/config"
	"gorm.io/driver/sqlite"
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

func Init() (err error) {
	dbFile := config.Get().DBFile
	if !strings.HasSuffix(dbFile, ".db") {
		return errors.New("無效的 db_file 名稱")
	}

	db, err = gorm.Open(
		sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", dbFile)),
		&gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "[GORM]", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
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
