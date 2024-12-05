package mysqlx

import (
	"fmt"

	"github.com/jiny3/gopkg/filex"
	"github.com/jiny3/gopkg/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	// for old code
	Mydb *gorm.DB
)

type DBConfig struct {
	User     string
	Password string
	Ip       string
	Port     int
	DBName   string
}

func init() {
	var dbconfig DBConfig
	err := filex.ReadConfig("config", "db", &dbconfig)
	if err != nil {
		logx.All.Error("read db config failed:", err)
		return
	}
	DB, err = NewDB(dbconfig)
	if err != nil {
		logx.All.Error("db connect failed:", err)
		return
	}
	logx.All.Info("db connect success")
	DB.AutoMigrate(&User{}, &Article{}, &Comment{})
	Mydb = DB
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
