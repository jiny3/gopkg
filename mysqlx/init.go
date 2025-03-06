package mysqlx

import (
	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/logx"
	"gorm.io/gorm"
)

// Deprecated: This global var will be removed in a future version.
// Use DB instead
var Mydb *gorm.DB

// Deprecated: This global var will be removed in a future version.
func init() {
	var dbconfig DBConfig
	err := configx.Read("config/db.yaml", &dbconfig)
	if err != nil {
		logx.All.Error("read db config failed:", err)
		return
	}

	db, err := New(dbconfig)
	if err != nil {
		logx.All.Error("db connect failed:", err)
		return
	}
	logx.All.Info("db connect success")
	db.AutoMigrate(&User{}, &Article{}, &Comment{})
	Mydb = db
	DB = db
}

// Deprecated: This function will be removed in a future version.
// Use New(dbconf) instead
func NewDB(dbconf DBConfig) (*gorm.DB, error) {
	return New(dbconf)
}
