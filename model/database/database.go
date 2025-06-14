package database

import (
	"github.com/sayasurvey/golang/model/schema"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
	"fmt"
)

var Db *gorm.DB
var err error

func DbInit() {
	godotenv.Load(".env")
	
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	
	sslMode := "require"
	if DB_HOST == "db" || DB_HOST == "localhost" {
		sslMode = "disable"
	}
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, sslMode)

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("database connection faild", err)
		panic("failed to connect to database")
	}

	Db.AutoMigrate(&schema.User{}, &schema.Book{}, &schema.BorrowedBook{}, &schema.BorrowingWishList{}, &schema.InvalidatedToken{})
	fmt.Println("gorm db connect")
}
