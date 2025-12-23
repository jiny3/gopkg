package dbx

import (
	"gorm.io/gorm"
)

// Select selects rows of type T from the database using GORM with the provided options.
func Select[T any](db *gorm.DB, table *T, opts ...SelectOption) ([]T, error) {
	tx := db.Model(table)
	for _, opt := range opts {
		opt(tx)
	}
	var result []T
	err := tx.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Insert inserts one or more rows of type T into the database using GORM.
func Insert[T any](db *gorm.DB, rows ...*T) error {
	if len(rows) == 0 {
		return nil
	}
	return db.Create(rows).Error
}

// Update updates rows of type T in the database using GORM with the provided options.
func Update[T any](db *gorm.DB, table *T, updates map[string]any, opts ...SelectOption) error {
	tx := db.Model(table)
	for _, opt := range opts {
		opt(tx)
	}
	return tx.Updates(updates).Error
}

// Delete deletes rows of type T from the database using GORM with the provided options.
func Delete[T any](db *gorm.DB, table *T, opts ...SelectOption) error {
	tx := db.Model(table)
	for _, opt := range opts {
		opt(tx)
	}
	return tx.Delete(table).Error
}
