package database

import (
	"github.com/varshard/mtl/infrastructure/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(conf *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DBConfig.ConnString))
	if err != nil {
		return db, err
	}

	return db, nil
}
