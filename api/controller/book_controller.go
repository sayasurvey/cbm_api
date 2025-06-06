package controller

import (
	"github.com/sayasurvey/golang/api/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sayasurvey/golang/model/schema"
	"github.com/sayasurvey/golang/model/database"
	"github.com/sayasurvey/golang/api/helper"
	"strconv"
)

type BookResponse struct {
	ID    		uint    	`json:"id"`
	Title  		string 		`json:"title"`
	ImageUrl 	string 		`json:"imageUrl"`
	Loanable 	bool 		`json:"loanable"`
	IsWishList  bool        `json:"isWishList"`
	User    struct {
		ID   uint    	`json:"id"`
		Name string 	`json:"name"`
	} `json:"user"`
}

type BooksResponse struct {
	Books      []BookResponse `json:"books"`
	CurrentPage int           `json:"currentPage"`
	LastPage    int           `json:"lastPage"`
	PerPage     int           `json:"perPage"`
}

func GetBooks(c *gin.Context) {
	page := 1
	perPage := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if perPageStr := c.Query("perPage"); perPageStr != "" {
		if parsedPerPage, err := strconv.Atoi(perPageStr); err == nil && parsedPerPage > 0 {
			perPage = parsedPerPage
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	books, err := repository.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の一覧取得に失敗しました",
		})
		return
	}

	wishListRepo := repository.NewBorrowingWishListRepository()
	wishList, err := wishListRepo.GetWishListByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "お気に入り情報の取得に失敗しました",
		})
		return
	}

	wishListMap := make(map[uint]bool)
	for _, item := range wishList {
		wishListMap[item.BookID] = true
	}

	var responseBooks []BookResponse
	for _, book := range books {
		responseUser := BookResponse{
			ID:       book.ID,
			Title:    book.Title,
			ImageUrl: book.ImageUrl,
			Loanable: book.Loanable,
			IsWishList: wishListMap[book.ID],
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

	paginatedBooks, currentPage, lastPage := helper.Pagination(responseBooks, page, perPage)

	c.JSON(http.StatusOK, BooksResponse{
		Books:      paginatedBooks,
		CurrentPage: currentPage,
		LastPage:    lastPage,
		PerPage:     perPage,
	})
}

type CreateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	ImageUrl string `json:"imageUrl" binding:"required"`
	Loanable bool   `json:"loanable"`
}

func CreateBook(c *gin.Context) {
	var request CreateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストボディが不正です",
		})
		return
	}

	var user schema.User
	if err := database.Db.First(&user, 1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー情報の取得に失敗しました",
		})
		return
	}

	book := schema.Book{
		UserId:   user.ID,
		Title:    request.Title,
		ImageUrl: request.ImageUrl,
		Loanable: request.Loanable,
		User:     user,
	}

	result := database.Db.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の作成に失敗しました",
		})
		return
	}

	response := BookResponse{
		ID:       book.ID,
		Title:    book.Title,
		ImageUrl: book.ImageUrl,
		Loanable: book.Loanable,
		User: struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		},
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "本の作成に成功しました",
		"book": response,
	})
}

type UpdateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	ImageUrl string `json:"imageUrl" binding:"required"`
	Loanable bool   `json:"loanable"`
}

type UpdateBookResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"imageUrl"`
	Loanable bool   `json:"loanable"`
	User     struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var request UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無効なリクエストです",
		})
		return
	}

	var book schema.Book
	if err := database.Db.Preload("User").First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	book.Title = request.Title
	book.ImageUrl = request.ImageUrl
	book.Loanable = request.Loanable

	if err := database.Db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の更新に失敗しました",
		})
		return
	}

	response := UpdateBookResponse{
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

	c.JSON(http.StatusOK, gin.H{
		"message": "本の更新に成功しました",
		"book": response,
	})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book schema.Book
	if err := database.Db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	if err := database.Db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の削除に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "本の削除に成功しました",
	})
}
