package mysqlx

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string
}

type Article struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
}

type Comment struct {
	gorm.Model
	Content   string
	UserID    uint
	ArticleID uint
}
