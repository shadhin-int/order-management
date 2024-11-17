package database

import (
	"gorm.io/gorm"
	"order-management/models"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Order{},
	)

}
