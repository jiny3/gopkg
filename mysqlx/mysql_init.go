package mysqlx

import (
	"fmt"

	"gitee.com/jiny1419/ucasnj-bbs/pkg/filex"
	"gitee.com/jiny1419/ucasnj-bbs/pkg/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MyDB *gorm.DB

func init() {
	var dbconfig DBConfig
	err := filex.ReadConfig("config", "db", &dbconfig)
	if err != nil {
		logx.MyAll.Error("read db config failed:", err)
		return
	}
	MyDB, err = NewDB(dbconfig)
	if err != nil {
		logx.MyAll.Error("db connect failed:", err)
		return
	}
	logx.MyAll.Info("db connect success")
	MyDB.AutoMigrate(&User{}, &Article{}, &Comment{})
}

func NewDB(dbconf DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbconf.User, dbconf.Password, dbconf.Ip, dbconf.Port, dbconf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
