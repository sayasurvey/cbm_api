package controller

import (
	"github.com/sayasurvey/golang/api/repository"
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"github.com/sayasurvey/golang/model/schema"
	"github.com/sayasurvey/golang/model/database"
)

type BookResponse struct {
	ID    		uint    `json:"id"`
	Title  		string 	`json:"title"`
	ImageUrl 	string 	`json:"image_url"`
	Loanable 	bool 		`json:"loanable"`
	User    struct {
		ID   uint    `json:"id"`
		Name string `json:"name"`
	} `json:"User"`
}

func GetBooks(context *gin.Context) {
	books, err := repository.GetAllBooks()
	fmt.Println("books", err)
	var responseBooks []BookResponse
	for _, book := range books {
		responseUser := BookResponse{
				ID:       book.ID,
				Title:    book.Title,
				ImageUrl: book.ImageUrl,
				Loanable: book.Loanable,
				User: struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				}{
					ID:   book.User.ID,
					Name: book.User.Name,
				},
		}
		responseBooks = append(responseBooks, responseUser)
	}
	context.JSON(http.StatusOK, responseBooks)
}

type CreateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	ImageUrl string `json:"image_url" binding:"required"`
	Loanable bool   `json:"loanable" binding:"required"`
}

func CreateBook(c *gin.Context) {
	var request CreateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	book := schema.Book{
		UserId:   1,
		Title:    request.Title,
		ImageUrl: request.ImageUrl,
		Loanable: request.Loanable,
	}

	result := database.Db.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create book",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"book": book,
	})
}
