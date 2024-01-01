package database

import (
	"github.com/varshard/mtl/infrastructure/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(conf *config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.ConnString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return db, err
	}

	return db, nil
}
