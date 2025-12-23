package dbx

import "gorm.io/gorm"

type option = func(*gorm.DB)

type optionReturnErr = func(*gorm.DB) error

type InitOption optionReturnErr

type SelectOption option

func WithAutoMigrate(models ...any) InitOption {
	return func(db *gorm.DB) error {
		return db.AutoMigrate(models...)
	}
}

func WithWhere(query any, args ...any) SelectOption {
	return func(db *gorm.DB) {
		db.Where(query, args...)
	}
}

func WithOrder(value string) SelectOption {
	return func(db *gorm.DB) {
		db.Order(value)
	}
}

func WithSelect(columns ...string) SelectOption {
	return func(db *gorm.DB) {
		db.Select(columns)
	}
}

func WithOmit(columns ...string) SelectOption {
	return func(db *gorm.DB) {
		db.Omit(columns...)
	}
}
