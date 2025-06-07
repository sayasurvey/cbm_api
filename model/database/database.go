package database

import (
	"github.com/sayasurvey/cbm_api/model/schema"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
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
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("database connection faild", err)
		panic("failed to connect to database")
	}

	Db.AutoMigrate(&schema.User{}, &schema.Book{}, &schema.BorrowedBook{}, &schema.BorrowingWishList{}, &schema.InvalidatedToken{})
	fmt.Println("gorm db connect")
}
