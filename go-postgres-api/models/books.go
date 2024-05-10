package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primary key;autoincrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

// With postgres, the database has to be created before using which can be done using automigrate function

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
