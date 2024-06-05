package dao

import (
	"time"
	"web/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	DB, err = gorm.Open(mysql.Open(config.MysqDSN), &gorm.Config{})
	if err != nil {
		return
	}

	if DB.Error != nil {
		return
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(1024)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
