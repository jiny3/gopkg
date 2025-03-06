package mysqlx

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// use it after New()
//
// TODO: add reconnect if connection lost
var DB *gorm.DB

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DbName   string
	Models   []any
}

func New(conf DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port)
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	sqlDB, err := _db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", conf.DbName))
	if err != nil {
		return nil, err
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(conf.Models...)
	return db, nil
}
