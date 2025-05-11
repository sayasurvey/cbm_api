package repository

import (
	"github.com/sayasurvey/golang/model/database"
	"github.com/sayasurvey/golang/model/schema"
	"fmt"
)

func GetAllBooks() ([]schema.Book, error) {
	var books []schema.Book
	db := database.Db

	result := db.Preload("User").Find(&books)
	if result.Error != nil {
		fmt.Println("エラー発生:", result.Error)
		return nil, result.Error
	}

	fmt.Printf("取得したユーザー数: %d\n", result.RowsAffected)
	return books, nil
}
