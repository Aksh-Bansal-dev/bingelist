package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	// db.Exec("DROP TABLE IF EXISTS shows, users, upvotes")
	db.AutoMigrate(&User{}, &Show{}, &Upvote{})
	return db
}
