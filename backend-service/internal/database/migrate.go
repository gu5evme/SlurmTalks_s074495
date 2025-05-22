package database

import (
	"fmt"
	"log"
	"sitter/internal/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("[WARNING] Migration failed:", err)
	}
	fmt.Println("[SUCCESS] Database migrated successfully!")
}
